package ast

type ExprVisitor interface {
	VisitBinaryExpr(b *BinaryExpr) interface{}
	VisitGroupingExpr(g *GroupingExpr) interface{}
	VisitLiteralExpr(l *LiteralExpr) interface{}
	VisitUnaryExpr(u *UnaryExpr) interface{}
	VisitVariableExpr(v *VariableExpr) interface{}
	VisitAssignExpr(a *AssignExpr) interface{}
	VisitLogicalExpr(l *LogicalExpr) interface{}
}

type StmtVisitor interface {
	VisitExprStmt(e *ExprStmt) interface{}
	VisitPrintStmt(p *PrintStmt) interface{}
	VisitVarStmt(p *VarStmt) interface{}
	VisitBlockStmt(b *BlockStmt) interface{}
	VisitIfStmt(i *IfStmt) interface{}
	VisitWhileStmt(w *WhileStmt) interface{}
}
