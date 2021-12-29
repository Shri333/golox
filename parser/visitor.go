package parser

type ExprVisitor interface {
	VisitBinaryExpr(b *BinaryExpr) interface{}
	VisitGroupingExpr(g *GroupingExpr) interface{}
	VisitLiteralExpr(l *LiteralExpr) interface{}
	VisitUnaryExpr(u *UnaryExpr) interface{}
	VisitVariableExpr(v *VariableExpr) interface{}
	VisitAssignExpr(a *AssignExpr) interface{}
	VisitLogicalExpr(l *LogicalExpr) interface{}
	VisitCallExpr(c *CallExpr) interface{}
	VisitGetExpr(g *GetExpr) interface{}
	VisitSetExpr(s *SetExpr) interface{}
	VisitThisExpr(t *ThisExpr) interface{}
	VisitSuperExpr(s *SuperExpr) interface{}
}

type StmtVisitor interface {
	VisitExprStmt(e *ExprStmt) interface{}
	VisitPrintStmt(p *PrintStmt) interface{}
	VisitVarStmt(v *VarStmt) interface{}
	VisitBlockStmt(b *BlockStmt) interface{}
	VisitIfStmt(i *IfStmt) interface{}
	VisitWhileStmt(w *WhileStmt) interface{}
	VisitFunStmt(f *FunStmt) interface{}
	VisitReturnStmt(r *ReturnStmt) interface{}
	VisitClassStmt(c *ClassStmt) interface{}
}
