import { ChildTypeDescriptor, designerFieldsMap } from "./designer-field-map";


export type DesignerFieldGroup = {
    main: ChildTypeDescriptor,
    children?: Array<ChildTypeDescriptor>
}

export const designerEdgesGroups: Array<DesignerFieldGroup> = [
    {
        main: designerFieldsMap.Relation,
    }
]

export const designerFieldsGroups: Array<DesignerFieldGroup> = [
    {
        main: designerFieldsMap.ShortText,
        children: [
            designerFieldsMap.ShortText,
            designerFieldsMap.LongText,
            designerFieldsMap.RichText,
            designerFieldsMap.Email
        ]
    }, {
        main: designerFieldsMap.Integer,
        children: [
            designerFieldsMap.Integer,
            designerFieldsMap.Decimal,
            designerFieldsMap.Float,
        ]
    }, {
        main: designerFieldsMap.DateTime,
        children: [
            designerFieldsMap.DateTime,
            designerFieldsMap.Date,
            designerFieldsMap.Time
        ]
    }, {
        main: designerFieldsMap.Boolean
    }, {
        main: designerFieldsMap.SingleChoice,
        children: [
            designerFieldsMap.SingleChoice,
            designerFieldsMap.MultipleChoice,
        ]
    }, {
        main: designerFieldsMap.media,

    }, {
        main: designerFieldsMap.Json
    }
];



