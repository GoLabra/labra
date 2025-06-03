import dayjs from 'dayjs';
import { MutableRefObject, ReactNode, useMemo, useRef } from 'react';
import { ApiFieldTypes, FormFieldTypes } from '@/types/field-type-descriptor';
import { localeConfig } from "@/config/locale-config"
import CheckIcon from '@mui/icons-material/Check';
import { ColumnDef, useResponsivePin, useRowExpansionStore } from 'mosaic-data-table';
import { Edge, Field, RelationType } from '@/lib/apollo/graphql.entities';
import { Avatar, Button, Chip, Stack, Typography } from '@mui/material';
// import { stringAvatar } from '@/lib/utils/avatar';
import { stringToDate, stringToDateTime, stringToTime } from '@/core-features/dynamic-form/value-convertor';
import { InlineAvatar } from '@/shared/components/avatar';
import { useViewRelationStore } from '@/core-features/view-item/use-view-relation-store';

export type ColumnOptions = {
    hasSort?: boolean,
    width?: number,
}

interface UseDynamicGridColumnsProps {
    entityName: string,
    fields?: Field[],
    edges?: Edge[],
    displayFieldName?: string,
	openRelation: (entityName: string, edge: Edge, entryId: string) => void,
    showId: boolean
}
export const useDynamicGridColumns = ({
    entityName,
    fields = [],
    edges = [],
    displayFieldName,
    openRelation,
    showId
}: UseDynamicGridColumnsProps): ColumnDef[] => {

    const displayFieldPin = useResponsivePin({ pin: 'left', breakpoint: 'sm', direction: 'up' });

    const expansionStoreRef = useRef(openRelation);
    expansionStoreRef.current = openRelation;

    return useMemo(() => (
        [
            ...fields
                .filter(i => showId ? true : i.name != 'id')
                .filter(i => !i.private) // TODO: this should be hidden from BE
                .filter(i => i.type != 'RichText') // TODO: set RichText based on configuration
                .filter(i => i.type != 'Json') // TODO: set JSON based on configuration
                .map(i => fieldToColumn(i))
                .map(i => i.id === displayFieldName ? {
                    ...i,
                    pin: displayFieldPin,
                    highlight: true
                } : i),
            ...edges.map(i => edgeToColumn(entityName, i, expansionStoreRef)),
        ]
    ), [displayFieldPin, showId, displayFieldName, fields, edges]);
}


const fieldToColumn = (field: Field): ColumnDef<Field> => {

    // Add avatar for name column
    if (field.name == 'name' && field.type == 'ShortText') {
        return shortTextColumnDef(field.name, field.caption, (row: any) => {
            const cellValue = row[field.name];
            return (<Stack direction="row" gap={1} alignItems="center">
                <InlineAvatar name={cellValue} />
                {cellValue}
            </Stack>)
        }, {
            width: 180,
            hasSort: true
        });
    }

    switch (field.type as ApiFieldTypes) {
        case 'ShortText': return shortTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'LongText': return longTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'RichText': return richTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Integer': return integerColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'DateTime': return dateTimeColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Date': return dateColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Time': return timeColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Boolean': return booleanColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'MultipleChoice': return multiChoiceColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: false });
        default: return shortTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
    }
}

const edgeToColumn = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>): ColumnDef<Edge> => {

    switch (edge.relationType) {
        case RelationType.One:
        case RelationType.OneToOne:
        case RelationType.OneToMany:
            return oneColumnDef(entityName, edge, openRelation);
        case RelationType.Many:
        case RelationType.ManyToOne:
        case RelationType.ManyToMany:
            return manyColumnDef(entityName, edge, openRelation);
    }

    return {
        id: edge.name,
        header: '',
        cell: (row: any) => ''
    }
}

const shortTextColumnDef = (name: string, caption: string, render: (row: any) => ReactNode, options?: ColumnOptions): ColumnDef<any> => {

    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const longTextColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 300,
        hasSort: options?.hasSort ?? false
    };
}

const richTextColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 300,
        hasSort: options?.hasSort ?? false
    };
}

const integerColumnDef = (name: string, caption: string, render: (row: any) => number, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 150,
        hasSort: options?.hasSort ?? false
    };
}

const dateTimeColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return stringToDateTime(value)?.format(localeConfig.dateTime.displayFormat) ?? undefined;

        },
        width: options?.width ?? 210,
        hasSort: options?.hasSort ?? false
    };
}

const dateColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return stringToDate(value)?.format(localeConfig.date.displayFormat) ?? undefined;

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const timeColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return stringToTime(value)?.format(localeConfig.time.displayFormat) ?? undefined;

        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const booleanColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return (<CheckIcon />)

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const multiChoiceColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);

            if (!Array.isArray(value)) {
                return value;
            }

            return (<Stack direction="row" gap="3px">
                {value.map((i) => {
                    return (<Chip key={i} label={i} size="small" variant="outlined" />)
                })}
            </Stack>
            )

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const oneColumnDef = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>, options?: ColumnOptions): ColumnDef<any> => {

    return {
        id: edge.name,
        header: edge.caption,
        cell: (row: any) => {
            const entityValue = row[edge.name];
            if (!entityValue) {
                return '';
            }

            const displayValue = entityValue[edge.relatedEntity.displayField.name] || `Id: ${entityValue['id']}`;

            return <Button
                variant="text"
                sx={{
                    padding: 0
                }}
                onClick={() => {
					openRelation.current(entityName, edge, row.id);
                }}>
                <Typography
                    color="text.primary"
                    variant='body2'
                    sx={{ textDecoration: 'underline' }}>{displayValue}</Typography></Button>
        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const manyColumnDef = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: edge.name,
        header: edge.caption,
        cell: (row: any) => {

            return <Button
                variant="text"
                sx={{
                    padding: 0
                }}
                onClick={() => {

					openRelation.current(entityName, edge, row.id);

                }}>
                <Typography
                    color="text.primary"
                    variant='body2'
                    sx={{ textDecoration: 'underline' }}>VIEW</Typography>
            </Button>
        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

