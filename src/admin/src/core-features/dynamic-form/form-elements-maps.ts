import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import React, { ReactNode } from "react";
import { NumberFormField } from "@/core-features/dynamic-form/form-fields/NumberField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { DateTimeFormField } from "@/core-features/dynamic-form/form-fields/DateTimeField";
import { TimeFormField } from "@/core-features/dynamic-form/form-fields/TimeField";
import { TextLongFormField } from "@/core-features/dynamic-form/form-fields/TextLongField";
import { RichTextFormField } from "@/core-features/dynamic-form/form-fields/RichTextField";
import { DateFormField } from "@/core-features/dynamic-form/form-fields/DateField";
import { JSONFormField } from "@/core-features/dynamic-form/form-fields/JSONField";
import { SelectFormField } from "@/core-features/dynamic-form/form-fields/SelectField";
import { FormFieldTypes } from "@/types/field-type-descriptor";
import { LookupOneFIELDFormField } from "@/core-features/dynamic-form/form-fields/LookupOneFIELD";
import { LookupManyFIELDFormField } from "./form-fields/LookupManyFIELD";
import { BooleanSelectFormField } from "@/core-features/dynamic-form/form-fields/BooleanSelectField";
import { SingleChoiceFormField } from "@/core-features/dynamic-form/form-fields/SingleChoice";
import { MultipleChoiceFormField } from "@/core-features/dynamic-form/form-fields/MultipleChoice";
import { TagsSelectFormField } from "@/core-features/dynamic-form/form-fields/TagsSelectField";
import { PasswordFormField } from "./form-fields/PasswordField";

export const FormElementsMap = { 
    ShortText: TextShortFormField,
    LongText: TextLongFormField,
    RichText: RichTextFormField,
    Password: PasswordFormField,
    
    Number: NumberFormField,

    Boolean: BooleanFormField,
    BooleanSelect: BooleanSelectFormField,

    DateTime: DateTimeFormField,
    Date: DateFormField,
    Time: TimeFormField,

    Json: JSONFormField,

    Select: SelectFormField,
    TagsSelect: TagsSelectFormField,

    SingleChoice: SingleChoiceFormField,
    MultipleChoice: MultipleChoiceFormField,
    RelationOne: LookupOneFIELDFormField,
    RelationMany: LookupManyFIELDFormField, 

} as {[key in FormFieldTypes]: React.FC<any>};