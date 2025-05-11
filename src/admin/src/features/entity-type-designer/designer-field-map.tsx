import NumbersIcon from "@/assets/icons/bootstrap/numbers";
import SwitchIcon from "@/assets/icons/bootstrap/switch";
import BracesIcon from "@/assets/icons/bootstrap/braces";
import { Edge, Field, RelationType } from "@/lib/apollo/graphql.entities";
import { PiSpiral } from "react-icons/pi";
// import { iconToNode } from "@/lib/utils/type-to-node";
import { TbRelationOneToManyFilled } from "react-icons/tb";
import { TbRelationOneToOneFilled } from "react-icons/tb";
import { PiGraphFill } from "react-icons/pi";
import { ApiFieldTypes, EdgeRelationTypes, FieldScalarTypes } from "@/types/field-type-descriptor";
import { MdShortText } from "react-icons/md";
import { LuText } from "react-icons/lu";
import { LuTextQuote } from "react-icons/lu";
import { MdAlternateEmail } from "react-icons/md";
import { PiMountainsFill } from "react-icons/pi";
import { TbClockFilled } from "react-icons/tb";
import { RiCalendarFill } from "react-icons/ri";
import { RiCalendarScheduleFill } from "react-icons/ri";
import { MdKey } from "react-icons/md";
import { GoMultiSelect } from "react-icons/go";
import { GoSingleSelect } from "react-icons/go";
import { MdPlaylistAddCheck } from "react-icons/md";
import { MdOutlineChecklistRtl } from "react-icons/md";
import { FC, ReactNode } from "react";
import DocumentIcon from "@/assets/icons/labra/document";
import { z } from "zod";

import { UseFormReturn } from "react-hook-form";
import { DesignerForm as shortTextDesignerForm, schema as shortTextSchema, schemaEffect as shortTextSchemaEffect, toDefaultValue as shortTextToDefaultValue } from "./designer-form/short-text-designer";
import { DesignerForm as longTextDesignerForm, schema as longTextSchema, schemaEffect as longTextSchemaEffect, toDefaultValue as longTextToDefaultValue } from "./designer-form/long-text-designer";
import { DesignerForm as rtfDesignerForm, schema as rtfSchema, schemaEffect as rtfSchemaEffect, toDefaultValue as rtfToDefaultValue } from "./designer-form/rtf-designer";
import { DesignerForm as emailDesignerForm, schema as emailSchema, schemaEffect as emailSchemaEffect, toDefaultValue as emailToDefaultValue } from "./designer-form/email-designer";
import { DesignerForm as numberDesignerForm, schema as numberSchema, schemaEffect as numberSchemaEffect, toDefaultValue as numberToDefaultValue } from "./designer-form/number-designer";
import { DesignerForm as dateTimeDesignerForm, schema as dateTimeSchema, schemaEffect as dateTimeSchemaEffect, toDefaultValue as dateTimeToDefaultValue } from "./designer-form/date-time-designer";
import { DesignerForm as dateDesignerForm, schema as dateSchema, schemaEffect as dateSchemaEffect, toDefaultValue as dateToDefaultValue } from "./designer-form/date-designer";
import { DesignerForm as timeDesignerForm, schema as timeSchema, schemaEffect as timeSchemaEffect, toDefaultValue as timeToDefaultValue } from "./designer-form/time-designer";
import { DesignerForm as booleanDesignerForm, schema as booleanSchema, schemaEffect as booleanSchemaEffect, toDefaultValue as booleanToDefaultValue } from "./designer-form/boolean-designer";
import { DesignerForm as jsonDesignerForm, schema as jsonSchema, schemaEffect as jsonSchemaEffect, toDefaultValue as jsonToDefaultValue } from "./designer-form/json-designer";
import { DesignerForm as singleChoiceDesignerForm, schema as singleChoiceSchema, schemaEffect as singleChoiceSchemaEffect, toDefaultValue as singleChoiceToDefaultValue } from "./designer-form/single-select-designer";
import { DesignerForm as multipleChoiceDesignerForm, schema as multipleChoiceSchema, schemaEffect as multipleChoiceSchemaEffect, toDefaultValue as multipleChoiceToDefaultValue } from "./designer-form/multiple-select-designer";
import { DesignerForm as relationDesignerForm, schema as relationSchema, schemaEffect as relationEffect, toDefaultValue as relationToDefaultValue } from "./designer-form/relation-designer";

import { ChangedFullEntity } from "@/types/entity";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";

export interface DesignerFormProps {
    formMethods: UseFormReturn<any, any, undefined>;
    openMode?: FormOpenMode;
}

export type ChildTypeDescriptor = {
    type: ApiFieldTypes;
    icon: JSX.Element;
    label: string;
    description?: string;
    disabled?: boolean,
    schema?: z.ZodObject<any, any>,
    schemaEffect?: null | ((schema: z.ZodObject<any, any>, editId: string | undefined, entity: ChangedFullEntity) => undefined | null | z.ZodEffects<any, any>),
    toDefaultValue?: (value: any) => any,
    formContent?: React.FC<{ formMethods: UseFormReturn<any, any, undefined>, openMode?: FormOpenMode }>
}

// export const designerEdges: Array<ChildTypeDescriptor> = [{
//     type: 'Relation',
//     icon: <PiGraphFill size={18} />,
//     label: 'Relation',
//     description: 'Foreign key to another entity',
// }]

export const designerFields: Array<ChildTypeDescriptor> = [{
    type: 'ID',
    icon: <MdKey size={18} />,
    label: 'Text Short',
    description: 'Enter brief text',
}, {
    type: 'ShortText',
    icon: <MdShortText size={18} />,
    label: 'Text Short',
    description: 'Enter brief text',
    //----
    schema: shortTextSchema,
    schemaEffect: shortTextSchemaEffect,
    toDefaultValue: shortTextToDefaultValue,
    formContent: shortTextDesignerForm
}, {
    type: 'LongText',
    icon: <LuText size={18} />,
    label: 'Text Long',
    description: 'Enter extended text',
    //----
    schema: longTextSchema,
    schemaEffect: longTextSchemaEffect,
    toDefaultValue: longTextToDefaultValue,
    formContent: longTextDesignerForm
}, {
    type: 'Email',
    icon: <MdAlternateEmail size={18} />,
    label: 'Email',
    description: 'Input email address',
    //----
    schema: emailSchema,
    schemaEffect: emailSchemaEffect,
    toDefaultValue: emailToDefaultValue,
    formContent: emailDesignerForm
}, {
    type: 'RichText',
    icon: <DocumentIcon />,
    label: 'Rich Text',
    description: 'Format text richly',
    //----
    schema: rtfSchema,
    schemaEffect: rtfSchemaEffect,
    toDefaultValue: rtfToDefaultValue,
    formContent: rtfDesignerForm
}, {
    type: 'Integer',
    icon: <NumbersIcon />,
    label: 'Number',
    description: 'Input numerical number',
    //----
    schema: numberSchema,
    schemaEffect: numberSchemaEffect,
    toDefaultValue: numberToDefaultValue,
    formContent: numberDesignerForm
}, {
    type: 'Decimal',
    icon: <NumbersIcon />,
    label: 'Decimal',
    description: 'Input decimal number',
    //----
    schema: numberSchema,
    schemaEffect: numberSchemaEffect,
    toDefaultValue: numberToDefaultValue,
    formContent: numberDesignerForm
}, {
    type: 'Float',
    icon: <NumbersIcon />,
    label: 'Float',
    description: 'Input floating number',
    //----
    schema: numberSchema,
    schemaEffect: numberSchemaEffect,
    toDefaultValue: numberToDefaultValue,
    formContent: numberDesignerForm
}, {
    type: 'DateTime',
    icon: <RiCalendarScheduleFill size={18} />,
    label: 'Date Time',
    description: 'Select date and time',
    //----
    schema: dateTimeSchema,
    schemaEffect: dateTimeSchemaEffect,
    toDefaultValue: dateTimeToDefaultValue,
    formContent: dateTimeDesignerForm
}, {
    type: 'Date',
    icon: <RiCalendarFill size={18} />,
    label: 'Date',
    description: 'Select a date',
    //----
    schema: dateSchema,
    schemaEffect: dateSchemaEffect,
    toDefaultValue: dateToDefaultValue,
    formContent: dateDesignerForm
}, {
    type: 'Time',
    icon: <TbClockFilled size={18} />,
    label: 'Time',
    description: 'Select a time',
    //----
    schema: timeSchema,
    schemaEffect: timeSchemaEffect,
    toDefaultValue: timeToDefaultValue,
    formContent: timeDesignerForm
}, {
    type: 'media',
    icon: <PiMountainsFill size={18} />,
    label: 'Media',
    description: 'Upload image or video',
    disabled: true
}, {
    type: 'Boolean',
    icon: <SwitchIcon />,
    label: 'Boolean',
    description: 'Toggle true/false',
    //----
    schema: booleanSchema,
    schemaEffect: booleanSchemaEffect,
    toDefaultValue: booleanToDefaultValue,
    formContent: booleanDesignerForm
}, {
    type: 'Json',
    icon: <BracesIcon />,
    label: 'JSON',
    description: 'Add JSON data',
    //----
    schema: jsonSchema,
    schemaEffect: jsonSchemaEffect,
    toDefaultValue: jsonToDefaultValue,
    formContent: jsonDesignerForm
}, {
    type: 'SingleChoice',
    icon: <MdPlaylistAddCheck size={18} />,
    label: 'Single Choice',
    description: 'Choose an option',
    //----
    schema: singleChoiceSchema,
    schemaEffect: singleChoiceSchemaEffect,
    toDefaultValue: singleChoiceToDefaultValue,
    formContent: singleChoiceDesignerForm
}, {
    type: 'MultipleChoice',
    icon: <MdOutlineChecklistRtl size={18} />,
    label: 'Multiple Choice',
    description: 'Choose an option',
    //----
    schema: multipleChoiceSchema,
    schemaEffect: multipleChoiceSchemaEffect,
    toDefaultValue: multipleChoiceToDefaultValue,
    formContent: multipleChoiceDesignerForm
}, {
    type: 'Relation',
    icon: <PiGraphFill size={18} />,
    label: 'Relation',
    description: 'Foreign key to another entity',
    //----
    schema: relationSchema,
    schemaEffect: relationEffect,
    toDefaultValue: relationToDefaultValue,
    formContent: relationDesignerForm
}];

type DesignerFieldMap = {
    [key in ApiFieldTypes]: ChildTypeDescriptor;
};

export const designerFieldsMap: Record<ApiFieldTypes, ChildTypeDescriptor> = designerFields.reduce((map, field) => {
    map[field.type] = field;
    return map;
}, {} as DesignerFieldMap);

// export const designerEdgeMap: Record<'Relation', ChildTypeDescriptor> = {
//     Relation: designerEdges[0]
// }

export const getIconForEntityChild = (child: Field | Edge): ReactNode => {

    const childType: ApiFieldTypes = 'type' in child ? child.type as ApiFieldTypes : 'Relation';
    const field = designerFieldsMap[childType]
    if (field) {
        return field.icon;
    }
    return <PiSpiral size={18} />;
    //}

    // if (child.__typename === 'Edge') {
    //     const field = designerEdgeMap.Relation;
    //     if (field) {
    //         return field.icon;
    //     }
    //     return <PiSpiral size={18} />;
    // }

    return <PiSpiral size={18} />;
}

export const getChildTypeForEntityChild = (child: Field | Edge): string => {

    if (child.__typename === 'Field') {
        return child.type.toUpperCase();
    }

    if (child.__typename === 'Edge') {
        return child.relatedEntity.caption?.toUpperCase();
    }

    return '';
}