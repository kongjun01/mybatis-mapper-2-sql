package ast

import (
	"bytes"
	"encoding/xml"

	"github.com/kongjun01/mybatis-mapper-2-sql/sqlfmt"
)

type QueryNode struct {
	*ChildrenNode
	Id   string
	Type string
}

func NewQueryNode() *QueryNode {
	n := &QueryNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (s *QueryNode) Scan(start *xml.StartElement) error {
	s.Type = start.Name.Local
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			s.Id = attr.Value
		}
	}
	return nil
}

func (s *QueryNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	ctx.QueryType = s.Type
	for _, a := range s.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(data)
	}
	return sqlfmt.FormatSQL(buff.String()), nil
}
