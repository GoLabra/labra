// src/lib/apollo/typeDefs.ts
export const typeDefs = `#graphql
  type CentrifugoMessage {
    id: ID!
    type: String!
    content: String!
    timestamp: String!
  }

  type Query {
    messages: [CentrifugoMessage!]!
  }

  type Mutation {
    addMessage(message: String!): Boolean!
    clearMessages: Boolean!
  }
`;