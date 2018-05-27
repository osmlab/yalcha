package osm

// An Object represents a Node, Way, Relation, Changeset, Note or User only.
type Object interface {
	ObjectID() int64

	// private is so that **ID types don't implement this interface.
	private()
}

func (n *Node) private()     {}
func (w *Way) private()      {}
func (r *Relation) private() {}

// Objects is a set of objects with some helpers
type Objects []Object
