import { localeConfig } from "@/config/locale-config";
import { stringToDate, stringToDateTime, stringToTime, timeToString } from "@/core-features/dynamic-form/value-convertor";
import { Counter, CounterGroup } from "@/shared/components/counter";
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

const getFlag = (show: boolean, flag: string, title: string) => {
    if (!show) {
        return null;
    }

    return (
        <Tooltip title="test">
            <Counter label={flag}></Counter>
        </Tooltip>
    );
};

interface EntityTypeDesignerFieldFlagsProps {
    child: DesignerField | DesignerEdge;
}
export const EntityTypeDesignerFieldFlags = (props: EntityTypeDesignerFieldFlagsProps) => {

    const {child} = props;


    const fieldFlags = useCallback((field: DesignerField) => (
        <CounterGroup size="10px">
            {getFlag(!!field.required, "R", "Required")}
            {getFlag(!!field.unique, "U", "Unique")}
            {getFlag(!!field.min, "M", `Min: ${formattedDefaultValue(field.type, field.min)}`)}
            {getFlag(!!field.max, "X", `Max: ${formattedDefaultValue(field.type, field.max)}`)}
            {getFlag(!!field.private, "P", "Private")}
            {getFlag(!!field.acceptedValues, "V", `Accepted Values: ${field.acceptedValues}`)}
            {getFlag(field.defaultValue !== undefined && field.defaultValue !== null, "D", `Default Value: ${formattedDefaultValue(field.type, field.defaultValue)}`)}
        </CounterGroup>
    ), []);

    const edgeFlags = useCallback((field: DesignerEdge) => (
        <CounterGroup size="10px">
            {getFlag(!!field.required, "R", "Required")}
            {getFlag(field.relationType === 'One', `➔ O`, "One")}
            {getFlag(field.relationType === 'Many', "➔ M", "Many")}
            {getFlag(field.relationType === 'OneToOne', "➔ OO", "One To One")}
            {getFlag(field.relationType === 'OneToMany', "➔ OM", "One To Many")}
            {getFlag(field.relationType === 'ManyToOne', "➔ MO", "Many To One")}
            {getFlag(field.relationType === 'ManyToMany', "➔ MM", "Many To Many")}

        </CounterGroup>), []);

    switch (child.__typename) {
        case "Field": return fieldFlags(child);
        case "Edge": return edgeFlags(child);
    }

    return null;
};
