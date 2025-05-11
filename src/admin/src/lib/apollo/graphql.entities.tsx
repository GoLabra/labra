export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateOnly: { input: any; output: any; }
  DateTime: { input: any; output: any; }
  TimeOnly: { input: any; output: any; }
};

export enum AppStatus {
  Fatal = 'FATAL',
  Generating = 'GENERATING',
  Restarting = 'RESTARTING',
  Reverting = 'REVERTING',
  Up = 'UP'
}

export type CreateEdgeInput = {
  belongsToCaption?: InputMaybe<Scalars['String']['input']>;
  caption: Scalars['String']['input'];
  private?: InputMaybe<Scalars['Boolean']['input']>;
  relatedEntity: EntityConnectInput;
  relationType: RelationType;
  required?: InputMaybe<Scalars['Boolean']['input']>;
};

export type CreateEntityInput = {
  caption: Scalars['String']['input'];
  displayField: FieldWhereUniqueInput;
  edges?: InputMaybe<CreateManyEdgesInput>;
  fields?: InputMaybe<CreateManyFieldsInput>;
};

export type CreateFieldInput = {
  acceptedValues?: InputMaybe<Array<InputMaybe<Scalars['String']['input']>>>;
  caption: Scalars['String']['input'];
  defaultValue?: InputMaybe<Scalars['String']['input']>;
  max?: InputMaybe<Scalars['String']['input']>;
  min?: InputMaybe<Scalars['String']['input']>;
  private?: InputMaybe<Scalars['Boolean']['input']>;
  required?: InputMaybe<Scalars['Boolean']['input']>;
  type: Scalars['String']['input'];
  unique?: InputMaybe<Scalars['Boolean']['input']>;
};

export type CreateManyEdgesInput = {
  create?: InputMaybe<Array<CreateEdgeInput>>;
};

export type CreateManyFieldsInput = {
  create?: InputMaybe<Array<CreateFieldInput>>;
};

export type Edge = {
  __typename?: 'Edge';
  belongsToCaption?: Maybe<Scalars['String']['output']>;
  caption: Scalars['String']['output'];
  name: Scalars['String']['output'];
  private?: Maybe<Scalars['Boolean']['output']>;
  relatedEntity: Entity;
  relationType?: Maybe<RelationType>;
  required?: Maybe<Scalars['Boolean']['output']>;
};

export type EdgeWhereUniqueInput = {
  caption?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type Entity = {
  __typename?: 'Entity';
  caption: Scalars['String']['output'];
  df: Scalars['String']['output'];
  displayField: Field;
  edges?: Maybe<Array<Edge>>;
  fields?: Maybe<Array<Field>>;
  name: Scalars['String']['output'];
  owner: EntityOwner;
  pluralName: Scalars['String']['output'];
};

export type EntityConnectInput = {
  connect: EntityWhereUniqueInput;
};

export enum EntityOwner {
  Admin = 'Admin',
  User = 'User'
}

export type EntityWhereUniqueInput = {
  caption?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type Field = {
  __typename?: 'Field';
  acceptedValues?: Maybe<Array<Maybe<Scalars['String']['output']>>>;
  caption: Scalars['String']['output'];
  defaultValue?: Maybe<Scalars['String']['output']>;
  max?: Maybe<Scalars['String']['output']>;
  min?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  private?: Maybe<Scalars['Boolean']['output']>;
  required?: Maybe<Scalars['Boolean']['output']>;
  type: Scalars['String']['output'];
  unique?: Maybe<Scalars['Boolean']['output']>;
};

export type FieldWhereUniqueInput = {
  caption?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  createEntity: Entity;
  deleteEntity: Entity;
  updateEntity: Entity;
};


export type MutationCreateEntityArgs = {
  data: CreateEntityInput;
};


export type MutationDeleteEntityArgs = {
  where: EntityWhereUniqueInput;
};


export type MutationUpdateEntityArgs = {
  data: UpdateEntityInput;
  where: EntityWhereUniqueInput;
};

export type Query = {
  __typename?: 'Query';
  entities?: Maybe<Array<Entity>>;
  entity?: Maybe<Entity>;
  fields?: Maybe<Array<Field>>;
};


export type QueryEntityArgs = {
  where?: InputMaybe<EntityWhereUniqueInput>;
};

export enum RelationType {
  Many = 'Many',
  ManyToMany = 'ManyToMany',
  ManyToOne = 'ManyToOne',
  One = 'One',
  OneToMany = 'OneToMany',
  OneToOne = 'OneToOne'
}

export type Subscription = {
  __typename?: 'Subscription';
  appStatus: AppStatus;
  entities?: Maybe<Array<Entity>>;
};

export type UpdateEdgeInput = {
  caption?: InputMaybe<Scalars['String']['input']>;
  private?: InputMaybe<Scalars['Boolean']['input']>;
  required?: InputMaybe<Scalars['Boolean']['input']>;
};

export type UpdateEntityInput = {
  caption?: InputMaybe<Scalars['String']['input']>;
  displayField?: InputMaybe<FieldWhereUniqueInput>;
  edges?: InputMaybe<UpdateManyEdgesInput>;
  fields?: InputMaybe<UpdateManyFieldsInput>;
};

export type UpdateFieldInput = {
  acceptedValues?: InputMaybe<Array<InputMaybe<Scalars['String']['input']>>>;
  caption?: InputMaybe<Scalars['String']['input']>;
  defaultValue?: InputMaybe<Scalars['String']['input']>;
  max?: InputMaybe<Scalars['String']['input']>;
  min?: InputMaybe<Scalars['String']['input']>;
  private?: InputMaybe<Scalars['Boolean']['input']>;
  required?: InputMaybe<Scalars['Boolean']['input']>;
  unique?: InputMaybe<Scalars['Boolean']['input']>;
};

export type UpdateManyEdgesInput = {
  create?: InputMaybe<Array<CreateEdgeInput>>;
  delete?: InputMaybe<Array<EdgeWhereUniqueInput>>;
  update?: InputMaybe<Array<UpdateOneEdgeInput>>;
};

export type UpdateManyFieldsInput = {
  create?: InputMaybe<Array<CreateFieldInput>>;
  delete?: InputMaybe<Array<FieldWhereUniqueInput>>;
  update?: InputMaybe<Array<UpdateOneFieldInput>>;
};

export type UpdateOneEdgeInput = {
  data: UpdateEdgeInput;
  where: EdgeWhereUniqueInput;
};

export type UpdateOneFieldInput = {
  data: UpdateFieldInput;
  where: FieldWhereUniqueInput;
};
