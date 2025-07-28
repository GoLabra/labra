import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { ChangedFullEntity } from "@/types/entity";
import { FormOpenMode, Options } from "@/core-features/dynamic-form/form-field";
import { SelectFormField } from "@/core-features/dynamic-form/form-fields/SelectField";

export const schema = z.object({
    caption: z.string().nonempty('Caption is required')
        .refine((val) => {
            const normalizedCaption = val.replace(/[\s_]/g, "").toLocaleLowerCase();
            return ENTITY_CHILDREN_SYSTEM_KEYWORDS.includes(normalizedCaption) == false;
        }, {
            message: 'System reserved keyword',
        })
        .refine((val) => {
            if (!val) {
                return true;
            }
            const firstLetter = val.charAt(0);
            return firstLetter.toLocaleLowerCase() != firstLetter.toLocaleUpperCase();
        }, {
            message: 'Caption must start with a letter',
        }),
	relationType: z.preprocess(i => i ?? '',
			z.string().min(1, `Relation Type is required`)
		),
    required: z.coerce.boolean()
});

export const toDefaultValue = (value: any): any => {
    return value;
}

export const schemaEffect = (schema: z.ZodObject<any, any>, editId: string | undefined, entity: ChangedFullEntity): z.ZodEffects<any, any> => {
    return schema.transform((data:any) => {

		return {
			...data,
			relatedEntity: {
				caption: 'File'
			}
		}
	}).refine((data) => {
        if (data.caption == undefined) {
            return false;
        }
        return checkChildCaptionExists(data.caption, editId, entity);
    }, {
        path: ['caption'],
        message: 'Caption already exists'
    })
}

const relationTypeOptions: Options = [{
    label: "One",
    value: "One"
}, {
    label: "Many",
    value: "Many"
}]

export const DesignerForm = (props: DesignerFormProps) => {

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
		<SelectFormField name="relationType" label="Relation Type" disabled={props.openMode == FormOpenMode.Edit} options={relationTypeOptions} required />
		<BooleanFormField name="required" label="Required" />
    </Stack>
    )
}
