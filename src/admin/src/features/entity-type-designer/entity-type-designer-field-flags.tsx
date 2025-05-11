import { localeConfig } from "@/config/locale-config";
import { stringToDate, stringToDateTime, stringToTime, timeToString } from "@/core-features/dynamic-form/value-convertor";
import { DesignerEdge, DesignerField } from "@/types/entity";
import { FieldScalarTypes } from "@/types/field-type-descriptor";
import { Box, Stack, styled, Tooltip } from "@mui/material";
import dayjs from "dayjs";
import { useCallback, useMemo } from "react";


const formattedDefaultValue = (type: string, value: any) => {
    if(!value){
        return null;
    }

    switch (type) {
        case 'DateTime':
            if (value == "now()") {
                return value;
            }
            return stringToDateTime(value)?.format(localeConfig.dateTime.displayFormat) ?? undefined;
            
        case 'Date':
            if (value == "now()") {
                return value;
            }
            return stringToDate(value)?.format(localeConfig.date.displayFormat) ?? undefined;

        case 'Time':
            if (value == "now()") {
                return value;
            }
            
            return stringToTime(value)?.format(localeConfig.time.displayFormat) ?? undefined;

        case 'Json': return JSON.stringify(value);

        case 'MultipleChoice': 
            return JSON.parse(value).join(', ');
    }

    return value;
}


const FlagBox = styled(Box)(({ theme }) => ({
    backgroundColor: theme.palette.background.default,
    fontSize: "10px",
    minWidth: "20px",
    height: "16px",
    display: "flex",
    padding: '0 7px',
    justifyContent: "center",
    alignItems: "center",
    opacity: 0.7,
    cursor: "default",
    transition: "opacity 0.2s",
    textWrap: 'nowrap',
    whiteSpace: 'nowrap',
    
    "&:hover": {
        opacity: 1,
    },
}));

const getFlag = (show: boolean, flag: string, title: string) => {
    if (!show) {
        return null;
    }

    return (
        <Tooltip title={title}>
            <FlagBox>{flag}</FlagBox>
        </Tooltip>
    );
};

interface EntityTypeDesignerFieldFlagsProps {
    child: DesignerField | DesignerEdge;
}
export const EntityTypeDesignerFieldFlags = (props: EntityTypeDesignerFieldFlagsProps) => {

    const {child} = props;


    const fieldFlags = useCallback((field: DesignerField) => (
        <Stack gap="1px" direction="row" borderRadius={1} overflow="hidden" alignItems="center">
            {getFlag(!!field.required, "R", "Required")}
            {getFlag(!!field.unique, "U", "Unique")}
            {getFlag(!!field.min, "M", `Min: ${formattedDefaultValue(field.type, field.min)}`)}
            {getFlag(!!field.max, "X", `Max: ${formattedDefaultValue(field.type, field.max)}`)}
            {getFlag(!!field.private, "P", "Private")}
            {getFlag(!!field.acceptedValues, "V", `Accepted Values: ${field.acceptedValues}`)}
            {getFlag(field.defaultValue !== undefined && field.defaultValue !== null, "D", `Default Value: ${formattedDefaultValue(field.type, field.defaultValue)}`)}
        </Stack>
    ), []);

    const edgeFlags = useCallback((field: DesignerEdge) => (
        <Stack gap="1px" direction="row" borderRadius={1} overflow="hidden" alignItems="center">
            {getFlag(!!field.required, "R", "Required")}
            {getFlag(field.relationType === 'One', `➔ O`, "One")}
            {getFlag(field.relationType === 'Many', "➔ M", "Many")}
            {getFlag(field.relationType === 'OneToOne', "➔ OO", "One To One")}
            {getFlag(field.relationType === 'OneToMany', "➔ OM", "One To Many")}
            {getFlag(field.relationType === 'ManyToOne', "➔ MO", "Many To One")}
            {getFlag(field.relationType === 'ManyToMany', "➔ MM", "Many To Many")}

        </Stack>), []);

    switch (child.__typename) {
        case "Field": return fieldFlags(child);
        case "Edge": return edgeFlags(child);
    }

    return null;
};
