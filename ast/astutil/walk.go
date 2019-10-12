// +build !appengine

package astutil

import (
	"fmt"
	"reflect"

	"github.com/gbl08ma/anko/ast"
)

// WalkFunc is used in Walk to walk the AST
type WalkFunc func(interface{}) error

// Walk walks the ASTs associated with a statement list generated by parser.ParseSrc
// each expression and/or statement is passed to the WalkFunc function.
// If the WalkFunc returns an error the walk is aborted and the error is returned
func Walk(stmt ast.Stmt, f WalkFunc) error {
	return walkStmt(stmt, f)
}

func walkStmts(stmts []ast.Stmt, f WalkFunc) error {
	for _, stmt := range stmts {
		if err := walkStmt(stmt, f); err != nil {
			return err
		}
	}
	return nil
}

func walkExprs(exprs []ast.Expr, f WalkFunc) error {
	for _, exp := range exprs {
		if err := walkExpr(exp, f); err != nil {
			return err
		}
	}
	return nil
}

func walkStmt(stmt ast.Stmt, f WalkFunc) error {
	//short circuit out if there are no functions
	if stmt == nil || f == nil {
		return nil
	}
	if err := callFunc(stmt, f); err != nil {
		return err
	}
	switch stmt := stmt.(type) {
	case *ast.StmtsStmt:
		if err := walkStmts(stmt.Stmts, f); err != nil {
			return err
		}
	case *ast.BreakStmt:
	case *ast.ContinueStmt:
	case *ast.LetMapItemStmt:
		if err := walkExpr(stmt.RHS, f); err != nil {
			return err
		}
		return walkExprs(stmt.LHSS, f)
	case *ast.ReturnStmt:
		return walkExprs(stmt.Exprs, f)
	case *ast.ExprStmt:
		return walkExpr(stmt.Expr, f)
	case *ast.VarStmt:
		return walkExprs(stmt.Exprs, f)
	case *ast.LetsStmt:
		if err := walkExprs(stmt.RHSS, f); err != nil {
			return err
		}
		return walkExprs(stmt.LHSS, f)
	case *ast.IfStmt:
		if err := walkExpr(stmt.If, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Then, f); err != nil {
			return err
		}
		if err := walkStmts(stmt.ElseIf, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Else, f); err != nil {
			return err
		}
	case *ast.TryStmt:
		if err := walkStmt(stmt.Try, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Catch, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Finally, f); err != nil {
			return err
		}
	case *ast.LoopStmt:
		if err := walkExpr(stmt.Expr, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Stmt, f); err != nil {
			return err
		}
	case *ast.ForStmt:
		if err := walkExpr(stmt.Value, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Stmt, f); err != nil {
			return err
		}
	case *ast.CForStmt:
		if err := walkStmt(stmt.Stmt1, f); err != nil {
			return err
		}
		if err := walkExpr(stmt.Expr2, f); err != nil {
			return err
		}
		if err := walkExpr(stmt.Expr3, f); err != nil {
			return err
		}
		if err := walkStmt(stmt.Stmt, f); err != nil {
			return err
		}
	case *ast.ThrowStmt:
		if err := walkExpr(stmt.Expr, f); err != nil {
			return err
		}
	case *ast.ModuleStmt:
		if err := walkStmt(stmt.Stmt, f); err != nil {
			return err
		}
	case *ast.SwitchStmt:
		if err := walkExpr(stmt.Expr, f); err != nil {
			return err
		}
		for _, switchCaseStmt := range stmt.Cases {
			caseStmt := switchCaseStmt.(*ast.SwitchCaseStmt)
			if err := walkStmt(caseStmt.Stmt, f); err != nil {
				return err
			}
		}
		if err := walkStmt(stmt.Default, f); err != nil {
			return err
		}
	case *ast.GoroutineStmt:
		return walkExpr(stmt.Expr, f)
	default:
		return fmt.Errorf("unknown statement %v", reflect.TypeOf(stmt))
	}
	return nil
}

func walkExpr(expr ast.Expr, f WalkFunc) error {
	//short circuit out if there are no functions
	if expr == nil || f == nil {
		return nil
	}
	if err := callFunc(expr, f); err != nil {
		return err
	}
	switch expr := expr.(type) {
	case *ast.OpExpr:
		return walkOperator(expr.Op, f)
	case *ast.LenExpr:
	case *ast.LiteralExpr:
	case *ast.IdentExpr:
	case *ast.MemberExpr:
		return walkExpr(expr.Expr, f)
	case *ast.ItemExpr:
		if err := walkExpr(expr.Item, f); err != nil {
			return err
		}
		return walkExpr(expr.Index, f)
	case *ast.SliceExpr:
		if err := walkExpr(expr.Item, f); err != nil {
			return err
		}
		if err := walkExpr(expr.Begin, f); err != nil {
			return err
		}
		return walkExpr(expr.End, f)
	case *ast.ArrayExpr:
		return walkExprs(expr.Exprs, f)
	case *ast.MapExpr:
		for i := range expr.Keys {
			if err := walkExpr(expr.Keys[i], f); err != nil {
				return err
			}
			if err := walkExpr(expr.Values[i], f); err != nil {
				return err
			}
		}
	case *ast.DerefExpr:
		return walkExpr(expr.Expr, f)
	case *ast.AddrExpr:
		return walkExpr(expr.Expr, f)
	case *ast.UnaryExpr:
		return walkExpr(expr.Expr, f)
	case *ast.ParenExpr:
		return walkExpr(expr.SubExpr, f)
	case *ast.FuncExpr:
		return walkStmt(expr.Stmt, f)
	case *ast.LetsExpr:
		if err := walkExprs(expr.LHSS, f); err != nil {
			return err
		}
		return walkExprs(expr.RHSS, f)
	case *ast.AnonCallExpr:
		if err := walkExpr(expr.Expr, f); err != nil {
			return err
		}
		return walkExpr(&ast.CallExpr{Func: reflect.Value{}, SubExprs: expr.SubExprs, VarArg: expr.VarArg, Go: expr.Go}, f)
	case *ast.CallExpr:
		return walkExprs(expr.SubExprs, f)
	case *ast.TernaryOpExpr:
		if err := walkExpr(expr.Expr, f); err != nil {
			return err
		}
		if err := walkExpr(expr.LHS, f); err != nil {
			return err
		}
		return walkExpr(expr.RHS, f)
	case *ast.ImportExpr:
		return walkExpr(expr.Name, f)
	case *ast.MakeExpr:
		if err := walkExpr(expr.LenExpr, f); err != nil {
			return err
		}
		return walkExpr(expr.CapExpr, f)
	case *ast.ChanExpr:
		if err := walkExpr(expr.RHS, f); err != nil {
			return err
		}
		return walkExpr(expr.LHS, f)
	case *ast.IncludeExpr:
		if err := walkExpr(expr.ItemExpr, f); err != nil {
			return err
		}
		return walkExpr(expr.ListExpr, f)
	default:
		return fmt.Errorf("unknown expression %v", reflect.TypeOf(expr))
	}
	return nil
}

func walkOperator(op ast.Operator, f WalkFunc) error {
	//short circuit out if there are no functions
	if op == nil || f == nil {
		return nil
	}
	if err := callFunc(op, f); err != nil {
		return err
	}
	switch op := op.(type) {
	case *ast.BinaryOperator:
		if err := walkExpr(op.LHS, f); err != nil {
			return err
		}
		return walkExpr(op.RHS, f)
	case *ast.ComparisonOperator:
		if err := walkExpr(op.LHS, f); err != nil {
			return err
		}
		return walkExpr(op.RHS, f)
	case *ast.AddOperator:
		if err := walkExpr(op.LHS, f); err != nil {
			return err
		}
		return walkExpr(op.RHS, f)
	case *ast.MultiplyOperator:
		if err := walkExpr(op.LHS, f); err != nil {
			return err
		}
		return walkExpr(op.RHS, f)
	}
	return nil
}

func callFunc(x interface{}, f WalkFunc) error {
	if x == nil || f == nil {
		return nil
	}
	return f(x)
}
