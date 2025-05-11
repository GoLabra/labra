import { FC, ReactNode } from "react";

export type FieldScalarTypes = 'ID' | 'ShortText' | 'LongText' | 'Email' | 'RichText' | 'Password' | 'Integer' | 'Decimal' | 'Float' | 'DateTime' | 'Date' | 'Time' | 'media' | 'Boolean' | 'Json' | 'Enum' | 'SingleChoice' | 'MultipleChoice';
export type EdgeRelationTypes = 'Relation';
export type ApiFieldTypes = FieldScalarTypes | EdgeRelationTypes;

export type FormFieldTypes = 'ShortText' | 'LongText' | 'RichText' | 'Password' | 'Number' | 'Boolean' | 'BooleanSelect' | 'DateTime' | 'Date' | 'Time' | 'Json' | 'Select' | 'TagsSelect' | 'SingleChoice' | 'MultipleChoice' | 'RelationOne' | 'RelationMany' | 'Select' | 'TagsSelect' | 'BooleanSelect' | 'Autocomplete' | 'EntitySelector' | 'RelationOne' | 'RelationMany';