package main

import "fmt"

type AstPrinter struct {
}

func (printer AstPrinter) Print(expr Expr) string {
	result := expr.Accept(printer)

	switch value := result.(type) {
	case string:
		return value
	default:
		return ""
	}
}

func (printer AstPrinter) VisitBinaryExpr(expr BinaryExpr) interface{} {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, expr.Left.Accept(printer), expr.Right.Accept(printer))
}

func (printer AstPrinter) VisitGroupingExpr(expr GroupingExpr) interface{} {
	return fmt.Sprintf("(group %s)", expr.Expression.Accept(printer))
}

func (AstPrinter) VisitLiteralExpr(expr LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}

	formatCode := ""

	switch expr.Value.(type) {
	case int:
		formatCode = "%d"
	case float32:
	case float64:
		formatCode = "%.2f"
	default:
		formatCode = "%s"
	}

	return fmt.Sprintf(formatCode, expr.Value)
}

func (printer AstPrinter) VisitUnaryExpr(expr UnaryExpr) interface{} {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, expr.Right.Accept(printer))
}
