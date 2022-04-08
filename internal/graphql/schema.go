package graphql

const graphqlSchema string = `

schema {
  query: Query
}
  
type Query {
  me: User!
  link(cursor: String!): Link!
  links(first: Int, last: Int, after: String): LinkConnection!
}
  
type User {
  id: ID!
  name: String!
  email: String!
  key: String!
}

type Link {
  cursor: String!
  title: String!
  url: String!
  byline: String!
  content: String!
  textcontent: String!
  excerpt: String!
  sitename: String!
  image: String
  favicon: String
  archive: [Int!]!
}

type LinkConnection {
  total: Int!
  nodes: [Link!]!
  edges: [LinkEdge!]!
  pageInfo: PageInfo!
}

type LinkEdge {
  cursor: String!
  node: Link!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

`
