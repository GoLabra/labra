"use client";

import type { FC } from 'react';
import { useCallback } from 'react';
import PropTypes from 'prop-types';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import { Button, Chip, Dialog, DialogActions, DialogContent, DialogTitle, Divider, IconButton, Stack, SvgIcon, Typography } from '@mui/material';
import { useFilters } from './use-filters';
import type { AdvancedFilter, AdvanedFilterOperator, AdvancedFilterProperty } from './filter';
import { FilterDialogItem } from './filter-dialog-item';

interface FilterDialogProps {
    filters?: AdvancedFilter[];
    onApply?: (filters: AdvancedFilter[]) => void;
    onClear?: () => void;
    onClose?: () => void;
    open?: boolean;
    properties?: AdvancedFilterProperty[];
}

export const FilterDialog: FC<FilterDialogProps> = (props) => {
    // Initial filters, we use the controlled state
    const { filters = [], onApply, onClear, onClose, open = false, properties = [] } = props;
    const filtersController = useFilters(properties, filters);

    const handleApply = useCallback(() => {
        // Allow the apply action only when every filter is considered valid.
        // Check the controller to configure the validation process.
        if (filtersController.valid) {
            onApply?.(filtersController.filters);
        }
    }, [filtersController.filters, filtersController.valid, onApply]);

    const displayClear = filters.length > 0;

    return (
        <Dialog
            onClose={onClose}
            open={open}
            PaperProps={{
                sx: {
                    maxWidth: 500,
                    width: '100%'
                }
            }}
        >
            <DialogTitle
                sx={{
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    display: 'flex'
                }}
            >
                <Typography variant="inherit">
                    Filter
                </Typography>
                <IconButton
                    onClick={onClose}
                    size="small"
                >
                    <SvgIcon fontSize="small">
                        <XMarkIcon />
                    </SvgIcon>
                </IconButton>
            </DialogTitle>
            <Divider />
            <DialogContent>
                <Stack
                    spacing={2}
                    divider={(
                        <Divider textAlign="left">
                            <Chip
                                color="primary"
                                label="AND"
                                size="small"
                            />
                        </Divider>
                    )}
                >
                    {filtersController.filters.map((filter, index) => {
                        const displayAdd = filtersController.filters.length === index + 1;

                        return (
                            <FilterDialogItem
                                displayAdd={displayAdd}
                                disableAdd={!filtersController.valid}
                                filter={filter}
                                index={index}
                                key={index}
                                onAdd={filtersController.handleFilterAdd}
                                onOperatorChange={filtersController.handleOperatorChange}
                                onPropertyChange={filtersController.handlePropertyChange}
                                onRemoveFilter={filtersController.handleFilterRemove}
                                onValueChange={filtersController.handleValueChange}
                                properties={properties}
                            />
                        );
                    })}
                </Stack>
            </DialogContent>
            <Divider />
            <DialogActions>
                {displayClear && (
                    <Button
                        color="error"
                        onClick={onClear}
                    >
                        Clear ({filters.length})
                    </Button>
                )}
                <Button
                    disabled={!filtersController.valid}
                    onClick={handleApply}
                    variant="contained"
                >
                    Apply
                </Button>
            </DialogActions>
        </Dialog>
    );
};

FilterDialog.propTypes = {
    onApply: PropTypes.func,
    onClear: PropTypes.func,
    onClose: PropTypes.func,
    open: PropTypes.bool,
    //operators: PropTypes.array,
    properties: PropTypes.array
};
