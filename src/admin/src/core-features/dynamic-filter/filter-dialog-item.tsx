"use client";

import type { FC } from 'react';
import { useMemo } from 'react';
import PropTypes from 'prop-types';
import PlusIcon from '@heroicons/react/24/outline/PlusIcon';
import { Button, FilledInput, MenuItem, Select, Stack, TextField, Typography } from '@mui/material';
// import { DatePicker } from '@mui/x-date-pickers';
import type { AdvancedFilter, AdvanedFilterOperator, AdvancedFilterProperty, AdvancedFilterValue } from './filter';

interface FilterDialogItemProps {
    disableAdd?: boolean;
    displayAdd?: boolean;
    filter: AdvancedFilter;
    index?: number;
    onAdd?: (index: number) => void;
    onOperatorChange?: (index: number, value: string) => void;
    onPropertyChange?: (index: number, name: string) => void;
    onRemoveFilter?: (index: number) => void;
    onValueChange?: (index: number, value: AdvancedFilterValue) => void;
    //operators?: FilterOperator[];
    properties?: AdvancedFilterProperty[];
}

export const FilterDialogItem: FC<FilterDialogItemProps> = (props) => {
    const {
        disableAdd = false,
        displayAdd = false,
        filter,
        index = 0,
        onAdd,
        onOperatorChange,
        onPropertyChange,
        onRemoveFilter,
        onValueChange,
        //operators = [],
        properties = []
    } = props;

    const property = useMemo(
        () => {
            return properties.find((property) => property.name === filter.property);
        },
        [filter, properties]
    );


    const operator = useMemo(
        () => {
            return property?.operators.find((operator) => operator.name === filter.operator);
        },
        [filter, property]
    );

    // const operatorOptions = useMemo(
    //     () => {
    //       return (property?.operators || [])
    //         .filter((operator) => !!operator) as FilterOperator[];
    //     },
    //     [property, property]
    //   );

    return (
        <Stack spacing={1}>
            <Typography variant="caption">
                Where
            </Typography>
            <Stack spacing={2}>
                <Stack
                    alignItems="center"
                    direction="row"
                    spacing={2}
                >
                    <Select
                        fullWidth
                        onChange={(event) => onPropertyChange?.(index, event.target.value)}
                        value={filter.property || ''}
                    >
                        {properties.map((property) => (
                            <MenuItem
                                key={property.name}
                                value={property.name}
                            >
                                {property.label}
                            </MenuItem>
                        ))}
                    </Select>
                    <Select
                        disabled={(property?.operators ?? []).length < 1}
                        fullWidth
                        displayEmpty
                        onChange={(event) => onOperatorChange?.(index, event.target.value)}
                        value={operator?.name || ''}
                    >
                        {property?.operators.map((operator) => (
                            <MenuItem
                                key={`item-${operator.name}`}
                                value={operator.name}
                            >
                                {operator.label}
                            </MenuItem>
                        ))}
                    </Select>
                </Stack>
                {/* {operator?.field === 'date' && (
          <DatePicker
            onChange={(date) => {
              if (date) {
                onValueChange?.(index, date);
              }
            }}
            renderInput={(inputProps) => (
              <TextField
                fullWidth
                {...inputProps}
              />
            )}
            value={filter.value || null}
          />
        )} */}
                {operator?.field === 'string' && (
                    <FilledInput
                        fullWidth
                        onChange={(event) => onValueChange?.(index, event.target.value)}
                        value={filter.value || ''}
                    />
                )}
                {operator?.field === 'number' && (
                    <FilledInput
                        fullWidth
                        type="number"
                        onChange={(event) => onValueChange?.(index, event.target.value)}
                        value={filter.value || ''}
                    />
                )}
                {operator?.field === 'boolean' && (
                     <Select
                        fullWidth
                        onChange={(event) => onValueChange?.(index, event.target.value)}
                        value={filter.value ?? true}
                    >
                        <MenuItem value={true as any}>True</MenuItem>
                        <MenuItem value={false as any}>False</MenuItem>
                    </Select>
                )}
            </Stack>
            <Stack
                alignItems="center"
                direction="row"
                spacing={2}
                justifyContent="flex-end"
            >
                {displayAdd && (
                    <Button
                        disabled={disableAdd}
                        onClick={() => onAdd?.(index + 1)}
                        size="small"
                        startIcon={<PlusIcon />}
                    >
                        Add
                    </Button>
                )}
                <Button
                    color="inherit"
                    onClick={() => onRemoveFilter?.(index)}
                    size="small"
                >
                    Remove
                </Button>
            </Stack>
        </Stack>
    );
};

FilterDialogItem.propTypes = {
    disableAdd: PropTypes.bool,
    displayAdd: PropTypes.bool,
    // @ts-ignore
    filter: PropTypes.object.isRequired,
    index: PropTypes.number,
    onAdd: PropTypes.func,
    onOperatorChange: PropTypes.func,
    onPropertyChange: PropTypes.func,
    onRemoveFilter: PropTypes.func,
    onValueChange: PropTypes.func,
    operators: PropTypes.array,
    properties: PropTypes.array
};
