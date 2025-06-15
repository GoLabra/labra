import { Button, DialogContent, IconButton, Stack } from "@mui/material";
import { SvgIcon } from "@mui/material";
import ArrowLeftIcon from "@heroicons/react/24/outline/ArrowLeftIcon";
import { forwardRef, useCallback, useContext, useEffect, useImperativeHandle, useMemo } from "react";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import EntityFieldTypeDescriptor from "./entity-field-type-descriptor";
import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { EntityTypeDesignerEntryDialogContent } from "./entity-type-designer-new-field-type-dialog";
import { Id } from "@/shared/components/id";
import { DesignerField } from "@/types/entity";
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { Form, useSchemaOptions } from "@/core-features/dynamic-form2/dynamic-form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from 'zod';
import { ChildTypeDescriptor, designerFieldsMap } from "./designer-field-map";
import { useSchemaEffectEntityParams } from "./use-designer-system-validation";

interface EntityTypeDesignerNewFieldPropsDialogProps {
    entityName: string;
    selectedChildTypeDescriptor: ChildTypeDescriptor;
    defaultValue: DesignerField;
}
export const EntityTypeDesignerNewFieldPropsDialog = forwardRef<ChainDialogContentRef, EntityTypeDesignerNewFieldPropsDialogProps>((props, ref) => {

    const { selectedChildTypeDescriptor, defaultValue } = props;
    const myDialogContext = useMyDialogContext();


    //const type = designerFieldsMap[props.selectedChildTypeDescriptor.type];
    if (!selectedChildTypeDescriptor) {
        throw new Error('Designer Form Type not found');
    }

    if (!selectedChildTypeDescriptor.schema) {
        throw new Error('Designer Form Type not found');
    }

    if (!selectedChildTypeDescriptor.formContent) {
        throw new Error('Designer Form Type not found');
    }

    const effectFunc = useSchemaEffectEntityParams(myDialogContext.editId, selectedChildTypeDescriptor.schemaEffect);
    const formDynamicSchema = useSchemaOptions(selectedChildTypeDescriptor.schema!, effectFunc);

    const formMethods = useForm({
        resolver: zodResolver(formDynamicSchema.pickedSchema),
        mode: 'all',
        defaultValues: selectedChildTypeDescriptor.toDefaultValue?.(defaultValue) ?? defaultValue
    });
    

    const onSetResult = useCallback(
        (data: any) => {
            myDialogContext.closeWithResults({
                selectedChildTypeDescriptor: selectedChildTypeDescriptor,
                field: {
                    ...data,
                    __typename: selectedChildTypeDescriptor.type == "Relation" ? "Edge" : "Field",
                    type: selectedChildTypeDescriptor.type,
                },
            });
        },
        [myDialogContext.closeWithResults],
    );

    const onSetResultAndAddAnother = useCallback(
        (data: any) => {
            onSetResult(data);
            myDialogContext.addPopupAsFirst(EntityTypeDesignerEntryDialogContent, {}, myDialogContext.openMode);
        },
        [onSetResult, myDialogContext.addPopupAsFirst],
    );

    const goToFieldType = useCallback(() => {
        myDialogContext.addPopupAsFirst(EntityTypeDesignerEntryDialogContent, {}, myDialogContext.openMode);
    }, [myDialogContext.addPopupAsFirst]);

    useImperativeHandle(ref, () => ({
        enterPressed: (e: React.KeyboardEvent<HTMLDivElement>) => {
            formMethods.handleSubmit(onSetResult)(e);
        }
    }));

    return (
        <>
            <DynamicDialogHeader>
                <Stack direction="row" gap={1}>
                    {myDialogContext.openMode === FormOpenMode.New && (
                        <IconButton onClick={goToFieldType}>
                            <SvgIcon fontSize="small">
                                <ArrowLeftIcon />
                            </SvgIcon>
                        </IconButton>
                    )}

                    {selectedChildTypeDescriptor && (
                        <EntityFieldTypeDescriptor
                            icon={selectedChildTypeDescriptor.icon}
                            label={selectedChildTypeDescriptor.label}
                        ></EntityFieldTypeDescriptor>
                    )}
                </Stack>
            </DynamicDialogHeader>

            <DialogContent>

                {myDialogContext.openMode != FormOpenMode.New && <Id id="Name" value={defaultValue.name} rootProps={{
                    marginLeft: 'auto',
                }} />}

                <Form methods={formMethods} dynamicSchema={formDynamicSchema} onSubmit={formMethods.handleSubmit(console.log)} >
                    <selectedChildTypeDescriptor.formContent formMethods={formMethods} openMode={myDialogContext.openMode} />
                </Form>


            </DialogContent>

            <DynamicDialogFooter>
                <Stack direction="row" gap={1}>
                    {myDialogContext.openMode === FormOpenMode.New && (
                        <Button
                            color="secondary"
                            variant="outlined"
                            onClick={formMethods.handleSubmit(onSetResultAndAddAnother)}
                        >
                            Save and Add Another Field
                        </Button>
                    )}

                    <Button
                        color="primary"
                        variant="contained"
                        onClick={formMethods.handleSubmit(onSetResult)}>
                        Save
                    </Button>
                </Stack>
            </DynamicDialogFooter>
        </>
    );
});
EntityTypeDesignerNewFieldPropsDialog.displayName = 'EntityTypeDesignerNewFieldPropsDialog';