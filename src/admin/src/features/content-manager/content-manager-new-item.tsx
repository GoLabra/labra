import { forwardRef, useCallback, useEffect, useImperativeHandle, useMemo } from "react";
import { Button, DialogContent, Stack, Typography } from "@mui/material";
import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { useFullEntity } from "@/hooks/use-entities";
import { useContentManagerFormSchema } from "./use-content-manager-form-schema";
import { Id } from "@/shared/components/id";
import { HistoryEvent, HistoryTimeLine } from "@/shared/components/history-timeline";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { localeConfig } from "@/config/locale-config";
import dayjs from 'dayjs';
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { InlineAvatar } from "@/shared/components/avatar";

interface ContentManagerEntryDialogContentProps {
    entityName: string,
    defaultValue: any
}
export const ContentManagerEntryDialogContent = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

    const myDialogContext = useMyDialogContext();
    const fullEntity = useFullEntity({ entityName: props.entityName });
    const formSchema = useContentManagerFormSchema(fullEntity);

    const id = useMemo(() => {
        return props.defaultValue?.id;
    }, [props.defaultValue]);

    const displayValue = useMemo(() => {
        if (!fullEntity) {
            return null;
        }

        const displayFieldProperty = fullEntity?.displayField?.name;
        if (!displayFieldProperty) {
            return null;
        }

        if (displayFieldProperty == 'id') {
            return null;
        }

        return props.defaultValue?.[displayFieldProperty];
    }, [fullEntity, props.defaultValue]);

    const upperValueAsFieldValues = useMemo(() => {
        return Object.fromEntries(
            Object.entries(myDialogContext.upperResults ?? {}).map(([key, value]) => {
                return [key, value];
            }))
    }, [myDialogContext.upperResults]);

    const defaultValue = useMemo(() => {
        return {
            ...formSchema.convertFromRawValue(formSchema.schemaDefaultValue),
            ...formSchema.convertFromRawValue(props.defaultValue),
            //...formSchema.convertFromRawValue(upperValueAsFieldValues
        };
        
    }, [props.defaultValue, upperValueAsFieldValues, formSchema.convertFromRawValue]);

    const formMethods = useForm({
        resolver: zodResolver(formSchema.schema),
        mode: 'all',
        defaultValues: defaultValue
    });

    const onSave = useCallback((formData: any) => {
        myDialogContext.closeWithResults(formData);
    }, [myDialogContext.closeWithResults]);

    const header = useMemo(() => {
        const fEntityName = fullEntity?.caption ?? props.entityName;
        if (!displayValue) {
            return fEntityName;
        }
        return `${fEntityName}: ${displayValue}`;
    }, [displayValue, fullEntity?.caption, props.entityName]);

    useImperativeHandle(ref, () => ({
        enterPressed: () => {
            formMethods.handleSubmit(onSave)();
        }
    }));

    return (
        <>
            <DynamicDialogHeader>
                {header}
            </DynamicDialogHeader>
            <DialogContent>
                <Stack
                    gap={2}>

                    {!!id && <Id value={id} rootProps={{
                        marginLeft: 'auto',
                    }} />}

                    <Form methods={formMethods} onSubmit={formMethods.handleSubmit(console.log)} >
                        <Stack gap={1.5} >
                            {formSchema.fields}
                        </Stack>
                    </Form>

                    {myDialogContext.openMode !== FormOpenMode.New && <BasicAuditTrail defaultValues={props.defaultValue} />}

                </Stack>
            </DialogContent>
            <DynamicDialogFooter>
                <Stack
                    direction="row"
                    gap={1}>
                    <Button
                        color="primary"
                        variant="contained"
                        onClick={formMethods.handleSubmit(onSave)}
                    >Save</Button>
                </Stack>
            </DynamicDialogFooter>
        </>
    )
});
ContentManagerEntryDialogContent.displayName = 'ContentManagerEntryDialogContent';

interface BasicAuditTrailProps {
    defaultValues: any;
}
const BasicAuditTrail = ({ defaultValues }: BasicAuditTrailProps) => {

    if (!defaultValues) {
        defaultValues = {};
    }

    const { createdBy, createdAt, updatedBy, updatedAt } = defaultValues;
    const timelineEvents = useMemo((): HistoryEvent[] => {
        return [{
            label: 'Last Updated',
            details: (
                <Stack>
                    <Typography
                        color="text.secondary"
                        fontSize={12}
                        lineHeight='1.5rem'>
                        {dayjs(updatedAt).format(localeConfig.dateTime.displayFormat)}
                    </Typography>

                    <Stack direction="row" gap={1}>
                        <InlineAvatar name={updatedBy?.name ?? 'Unknown'} />
                        <Typography color="text.secondary">{updatedBy?.name ?? 'Unknown'}</Typography>
                    </Stack>

                </Stack>)
        }, {
            label: 'Created',
            details: (
                <Stack>
                    <Typography
                        color="text.secondary"
                        fontSize={12}
                        lineHeight='1.5rem'>
                        {dayjs(createdAt).format(localeConfig.dateTime.displayFormat)}
                    </Typography>

                    <Stack direction="row" gap={1}>
                        <InlineAvatar name={createdBy?.name ?? 'Unknown'} />
                        <Typography color="text.secondary">{createdBy?.name ?? 'Unknown'}</Typography>
                    </Stack>
                    
                </Stack>)

        }];
    }, [createdBy, createdAt, updatedBy, updatedAt]);

    if (!defaultValues) {
        return null;
    }

    return (
        <HistoryTimeLine events={timelineEvents} />
    )
}