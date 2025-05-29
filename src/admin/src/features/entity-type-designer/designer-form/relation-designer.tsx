import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React, { useMemo } from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { TimeFormField } from "@/core-features/dynamic-form/form-fields/TimeField";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { ChangedFullEntity } from "@/types/entity";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { stringToTime, timeToString } from "@/core-features/dynamic-form/value-convertor";
import { SelectFormField } from "@/core-features/dynamic-form/form-fields/SelectField";
import { useEntitiesDesigner } from "../use-designer-entities";
import { FormOpenMode, Options } from "@/core-features/dynamic-form/form-field";

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
    relatedEntity: z.preprocess(i => {
        if (!i) {
            return { caption: '' };
        }
        return i;
    },
        z.object({
            caption: z.preprocess(i => i ?? '',
                z.string().min(1, `Related Entity is required`)
            )
        })
    ),
    belongsToCaption: z.preprocess(i => i ?? '',
				z.string().min(1, 'Belongs To Caption is required')
			),
    required: z.coerce.boolean(),
});

export const schemaEffect = (schema: z.ZodObject<any, any>, editId: string | undefined, entity: ChangedFullEntity): z.ZodEffects<any, any> => {
    return schema.refine((data) => {
        if (data.caption == undefined) {
            return false;
        }
        return checkChildCaptionExists(data.caption, editId, entity);
    }, {
        path: ['caption'],
        message: 'Caption already exists'
    })
}

export const toDefaultValue = (value: any): any => {
    if (!value) {
        return value;
    }

    return {
        ...value
    }
}


const relationTypeOptions: Options = [{
    label: "One",
    value: "One"
}, {
    label: "Many",
    value: "Many"
}, {
    label: "One To One",
    value: "OneToOne"
}, {
    label: "One To Many",
    value: "OneToMany"
}, {
    label: "Many To One",
    value: "ManyToOne"
}, {
    label: "Many To Many",
    value: "ManyToMany"
}]

const belongsToCaptionVisible = ['OneToOne', 'OneToMany', 'ManyToOne', 'ManyToMany',];

export const DesignerForm = (props: DesignerFormProps) => {

    const { allEntities } = useEntitiesDesigner();

    const options = useMemo((): Options => {
        return allEntities.map((entity) => ({
            label: entity.caption!,
            value: entity.caption! // we will use the caption to connect to the related entity
        }));
    }, [allEntities]);

    const relationType = props.formMethods.watch('relationType');
    const belongsToCaptionHide = useMemo(() => !belongsToCaptionVisible.includes(relationType), [relationType]);

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <SelectFormField name="relationType" label="Relation Type" disabled={props.openMode == FormOpenMode.Edit} options={relationTypeOptions} required />
        <Stack direction="row" gap={1.5}>
            <SelectFormField name="relatedEntity.caption" label="Related Entity" disabled={props.openMode == FormOpenMode.Edit} options={options} required />
            <TextShortFormField name="belongsToCaption" label="Belongs To Caption" disabled={props.openMode == FormOpenMode.Edit} hide={belongsToCaptionHide} required />
        </Stack>
        <BooleanFormField name="required" label="Required" />
    </Stack>
    )
}
