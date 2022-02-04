# FORK

NOTE: This internal package is a fork of https://github.com/shurcooL/graphql

Shurcool's excellent, MIT-licensed work was needlessly tied to an HTTP client,
causing some unnecessary complexity downstream. This fork simply strips the
HTTP support, and exposes the underlying marshalling/encoding machinery
directly.

Docs below have been revised, since the `Client` struct no longer exists. However,
the the reflective encoding stuff is _mostly_ unchanged. See the
[github.com/deref/graphql-go/encoding](https://github.com/deref/graphql-go/tree/main/encoding)
package for the public interface.

One critical change is improved scalar support. In particular, custom wrapper
types are no longer required for representing GraphQL's `String`, `Boolean`,
`Float`, and `Int` types. You can use corresponding Go types directly, provided
they can represent the same range of values. That is, any integer type that
fits in a 64-bit float can be used.  The machine-specific `int` type may be
truncated. Similar goes for floating point values. Other custom scalar types
support both encoding and decoding by implementing Go's standard json.Marshaler
and json.Unmarshaler interfaces.

### Simple Query

To make a GraphQL query, you need to define a corresponding Go type.

For example, to make the following GraphQL query:

```GraphQL
query {
	me {
		name
	}
}
```

You can define this variable:

```Go
var query struct {
	Me struct {
		Name string
	}
}
```

Then call `client.Query`, passing a pointer to it:

```Go
err := client.Query(context.Background(), &query, nil)
if err != nil {
	// Handle error.
}
fmt.Println(query.Me.Name)

// Output: Luke Skywalker
```

### Arguments and Variables

Often, you'll want to specify arguments on some fields. You can use the `graphql` struct field tag for this.

For example, to make the following GraphQL query:

```GraphQL
{
	human(id: "1000") {
		name
		height(unit: METER)
	}
}
```

You can define this variable:

```Go
var q struct {
	Human struct {
		Name   string
		Height float64 `graphql:"height(unit: METER)"`
	} `graphql:"human(id: \"1000\")"`
}
```

Then call `client.Query`:

```Go
err := client.Query(context.Background(), &q, nil)
if err != nil {
	// Handle error.
}
fmt.Println(q.Human.Name)
fmt.Println(q.Human.Height)

// Output:
// Luke Skywalker
// 1.72
```

However, that'll only work if the arguments are constant and known in advance. Otherwise, you will need to make use of variables. Replace the constants in the struct field tag with variable names:

```Go
var q struct {
	Human struct {
		Name   string
		Height float64 `graphql:"height(unit: $unit)"`
	} `graphql:"human(id: $id)"`
}
```

Then, define a `variables` map with their values:

```Go
variables := map[string]interface{}{
	"id":   graphql.ID(id),
	"unit": starwars.LengthUnit("METER"),
}
```

Finally, call `client.Query` providing `variables`:

```Go
err := client.Query(context.Background(), &q, variables)
if err != nil {
	// Handle error.
}
```

### Inline Fragments

Some GraphQL queries contain inline fragments. You can use the `graphql` struct field tag to express them.

For example, to make the following GraphQL query:

```GraphQL
{
	hero(episode: "JEDI") {
		name
		... on Droid {
			primaryFunction
		}
		... on Human {
			height
		}
	}
}
```

You can define this variable:

```Go
var q struct {
	Hero struct {
		Name  string
		Droid struct {
			PrimaryFunction string
		} `graphql:"... on Droid"`
		Human struct {
			Height float64
		} `graphql:"... on Human"`
	} `graphql:"hero(episode: \"JEDI\")"`
}
```

Alternatively, you can define the struct types corresponding to inline fragments, and use them as embedded fields in your query:

```Go
type (
	DroidFragment struct {
		PrimaryFunction string
	}
	HumanFragment struct {
		Height float64
	}
)

var q struct {
	Hero struct {
		Name          string
		DroidFragment `graphql:"... on Droid"`
		HumanFragment `graphql:"... on Human"`
	} `graphql:"hero(episode: \"JEDI\")"`
}
```

Then call `client.Query`:

```Go
err := client.Query(context.Background(), &q, nil)
if err != nil {
	// Handle error.
}
fmt.Println(q.Hero.Name)
fmt.Println(q.Hero.PrimaryFunction)
fmt.Println(q.Hero.Height)

// Output:
// R2-D2
// Astromech
// 0
```

### Mutations

Mutations often require information that you can only find out by performing a query first. Let's suppose you've already done that.

For example, to make the following GraphQL mutation:

```GraphQL
mutation($ep: Episode!, $review: ReviewInput!) {
	createReview(episode: $ep, review: $review) {
		stars
		commentary
	}
}
variables {
	"ep": "JEDI",
	"review": {
		"stars": 5,
		"commentary": "This is a great movie!"
	}
}
```

You can define:

```Go
var m struct {
	CreateReview struct {
		Stars      int
		Commentary string
	} `graphql:"createReview(episode: $ep, review: $review)"`
}
variables := map[string]interface{}{
	"ep": starwars.Episode("JEDI"),
	"review": starwars.ReviewInput{
		Stars:      5,
		Commentary: "This is a great movie!",
	},
}
```

Then call `client.Mutate`:

```Go
err := client.Mutate(context.Background(), &m, variables)
if err != nil {
	// Handle error.
}
fmt.Printf("Created a %v star review: %v\n", m.CreateReview.Stars, m.CreateReview.Commentary)

// Output:
// Created a 5 star review: This is a great movie!
```

## Directories

| Path                                                                                   | Synopsis                                                                                                        |
| -------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| [example/graphqldev](https://godoc.org/github.com/shurcooL/graphql/example/graphqldev) | graphqldev is a test program currently being used for developing graphql package.                               |
| [ident](https://godoc.org/github.com/shurcooL/graphql/ident)                           | Package ident provides functions for parsing and converting identifier names between various naming convention. |
| [internal/jsonutil](https://godoc.org/github.com/shurcooL/graphql/internal/jsonutil)   | Package jsonutil provides a function for decoding JSON into a GraphQL query data structure.                     |

## License

- [MIT License](LICENSE)
