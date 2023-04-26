//DO NOT EDIT THIS FILE
//IT WAS GENERATED BY CODE UNDER TOOLS FOLDER

package main

type Visitor interface {
	VisitBinaryExpr(expr BinaryExpr) interface{}
	VisitGroupingExpr(expr GroupingExpr) interface{}
	VisitLiteralExpr(expr LiteralExpr) interface{}
	VisitUnaryExpr(expr UnaryExpr) interface{}
}

type Expr interface {
	Accept(visitor Visitor) interface{}
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr BinaryExpr) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(expr)
}

type GroupingExpr struct {
	Expression Expr
}

func (expr GroupingExpr) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(expr)
}

type LiteralExpr struct {
	Value interface{}
}

func (expr LiteralExpr) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr UnaryExpr) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(expr)
}
