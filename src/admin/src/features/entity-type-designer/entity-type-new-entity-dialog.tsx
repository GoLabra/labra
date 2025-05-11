import { Box, Button, DialogContent, Stack, Typography } from "@mui/material";
import { forwardRef, useCallback, useImperativeHandle, useMemo } from "react";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { ChangedNameCaptionEntity } from "@/types/entity";
import { Id } from "@/shared/components/id";
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import Avvvatars from "avvvatars-react";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import pluralize from "pluralize";
import { ENTITY_SYSTEM_KEYWORDS } from "@/config/CONST";
import { useEntitiesDesigner } from "./use-designer-entities";

export const getSchema = (entitiesDesigner: ReturnType<typeof useEntitiesDesigner>, editId?: string) =>
    z.object({
        caption: z.string()
            .nonempty('Caption is required')
            .refine((val) => {
                const normalizedCaption = val.replace(/[\s_]/g, "").toLocaleLowerCase();
                const existingField = entitiesDesigner.allEntities?.find(i => i.caption?.toLocaleLowerCase() == normalizedCaption && i.name != editId);
                return !existingField;
            }, {
                message: 'Entity already exists'
            })
            .refine((val) => {
                const normalizedCaption = val.replace(/[\s_]/g, "").toLocaleLowerCase();
                return ENTITY_SYSTEM_KEYWORDS.includes(normalizedCaption) == false;
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
            })
            .refine((val) => pluralize(val) != val, {
                message: 'Caption must be singular',
            })
            .refine((val) => !val?.toLowerCase().endsWith('test'), {
                message: 'Caption cannot end with "test"',
            })


    });

interface EntityTypeNewEntityDialogProps {
    defaultValue: ChangedNameCaptionEntity;
}
export const EntityTypeNewEntityDialog = forwardRef<ChainDialogContentRef, EntityTypeNewEntityDialogProps>((props, ref) => {

    const { defaultValue } = props;
    const myDialogContext = useMyDialogContext();
    const entitiesDesigner = useEntitiesDesigner();

    const schema = useMemo(() => getSchema(entitiesDesigner, myDialogContext.editId), [entitiesDesigner, myDialogContext.editId]);

    const formMethods = useForm({
        resolver: zodResolver(schema),
        mode: 'all',
        defaultValues: defaultValue
    });

    const onFinish = useCallback((data: any) => {
        myDialogContext.closeWithResults({ ...data });
    }, [myDialogContext.closeWithResults]);

    useImperativeHandle(ref, () => ({
        enterPressed: () => {
            formMethods.handleSubmit(onFinish)();
        }
    }));

    return (<>

        <DynamicDialogHeader whatsThis="An entity represents a distinct piece of content or data within the CMS.
                It can be anything from a blog post, product, user profile, or any other type of data you need to manage.
                Entities help organize and structure your content efficiently.">
            {myDialogContext.openMode == FormOpenMode.New ? 'New Entity' : 'Edit Entity'}
        </DynamicDialogHeader>

        <DialogContent sx={{ overflow: 'visible' }}>

            <Stack direction="row" justifyContent="center" sx={{
                width: 'fit-content',
                margin: 'auto',
                padding: '5px',
                border: '2px dashed color-mix(in srgb, var(--mui-palette-background-paper), var(--mui-palette-common-onBackground) 20%)',
                borderRadius: '100%'
            }}>
                <Avvvatars value={formMethods.watch('caption') ?? 'Unknown'} style="shape" size={65} />
            </Stack>

            {defaultValue?.name && <Id id="Name" value={defaultValue.name} rootProps={{
                marginLeft: 'auto',
            }} />}

            <Form methods={formMethods} onSubmit={formMethods.handleSubmit(console.log)} >
                <TextShortFormField name="caption" label="Caption" required />
            </Form>

        </DialogContent>
        <DynamicDialogFooter>
            <Button
                color="primary"
                variant="contained"
                onClick={formMethods.handleSubmit(onFinish)}
            >Save</Button>
        </DynamicDialogFooter>
    </>)
});
EntityTypeNewEntityDialog.displayName = 'EntityTypeNewEntityDialog';