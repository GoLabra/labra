"use client";

import PropTypes from 'prop-types';
import FilterAltIcon from '@mui/icons-material/FilterAlt';
import ChecklistIcon from '@mui/icons-material/Checklist';
import CachedIcon from '@mui/icons-material/Cached';
import { Button, IconButton, ListItemIcon, ListItemText, MenuItem, Stack, SvgIcon, Tooltip } from '@mui/material';
import { FieldScalarTypes } from '@/types/field-type-descriptor';
import { AdvancedFilter as AdvancedFilter } from '@/core-features/dynamic-filter/filter';
import { MenuButton } from '@/shared/components/menu/menu-button';
import { MenuItemSelect } from '@/shared/components/menu/menu-item-select';
import { BulkActionsMenu } from '@/shared/components/bulk-actions-menu';
import { QueryField } from '@/shared/components/query-field';
import { FilterEditor } from '@/shared/components/filter-editor';
import { ColumnDef } from 'mosaic-data-table';
import WrenchScrewdriverIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';
import { MenuItemToggle } from '@/shared/components/menu/menu-item-toggle';
import { useContentManagerSearch } from '@/hooks/use-content-manager-search';
import { ResponsiveButton } from '@/styles/button.responsive';

export type SearchField = {
    label: string;
    name: string;
    type: FieldScalarTypes;
    selectOptions?: string[];
}

interface OrdersSearchProps {
    //fields: SearchField[];
    headCells?: ColumnDef[];
    disabled?: boolean;
    onRefresh?: () => void;

    filters?: AdvancedFilter[];
    onAdvancedFiltersApply?: (filters: AdvancedFilter[]) => void;
    onAdvancedFiltersClear?: () => void;

    onQueryChange?: (query: string) => void;
    query?: string;
    selected?: string[];
    onBulkDelete?: () => void;

    selectionEnabled: boolean;
    onSelectionEnabledChange?: (enabled: boolean) => void;

    filterEnabled: boolean;
    onFilterEnabledChange?: (enabled: boolean) => void;

    contentManagerSearch: ReturnType<typeof useContentManagerSearch>;

    showId: boolean;
    onShowIdChange: (showId: boolean) => void;
}
export const ContentManagerSearch = (props: OrdersSearchProps) => {

    const { headCells, disabled = false, filters = [], onQueryChange, query = '', selected = [], onBulkDelete } = props;
    const hasSelection = selected.length > 0;

    return (
        <>
            <div>
                <Stack gap={2}>

                    <Stack
                        alignItems="center"
                        direction="row"
                        flexWrap="wrap"
                        gap={2}
                    >
                        {hasSelection && (
                            <BulkActionsMenu
                                disabled={disabled}
                                selectedCount={selected.length}
                                onDelete={onBulkDelete}
								sx={{
									order: {
										xs: 1,
										sm: 0
									}
								}}
                            />
                        )}

                        <Tooltip title="Refresh" placement='bottom' arrow
							sx={{
								order: {
									xs: 2,
									sm: 1
								}
							}}>
                            <IconButton
                                color="primary"
                                onClick={() => props.onRefresh?.()}
                                size="medium">
                                <CachedIcon />
                            </IconButton>
                        </Tooltip>

                        <QueryField
                            disabled={disabled}
                            placeholder="Search..."
                            onChange={onQueryChange}
                            sx={{
                                flexGrow: 1,
								minWidth: {
									xs: '100%',
									sm: 'inherit'
								},
								order: {
									xs: 0,
									sm: 2
								}
                            }}
                            value={query}
                        />

                        <Stack
                            direction="row"
                            gap={2}
                            justifyContent="end"
							sx={{
								order: {
									xs: 3,
									sm: 3
								}
							}}>

                            <ResponsiveButton
                                disabled={disabled}
                                onClick={() => props.onSelectionEnabledChange?.(!props.selectionEnabled)}
                                size="medium"
                                startIcon={(
                                    <SvgIcon fontSize="small">
                                        <ChecklistIcon />
                                    </SvgIcon>
                                )}
                                variant={props.selectionEnabled ? 'contained' : 'text'}
                                aria-haspopup="dialog"
                            >
                                <span className="button-text">Selection</span>
                            </ResponsiveButton>

                            <ResponsiveButton
                                disabled={disabled}
                                onClick={() => {
                                    var newValue = !props.filterEnabled;
                                    props.onFilterEnabledChange?.(newValue);
                                    if (!newValue) {
                                        props.contentManagerSearch.handleFiltersApply({});
                                    }
                                }}
                                size="medium"
                                startIcon={(
                                    <SvgIcon fontSize="small">
                                        <FilterAltIcon />
                                    </SvgIcon>
                                )}
                                variant={props.filterEnabled ? 'contained' : 'text'}
                                aria-haspopup="dialog"
                            >
                                <span className="button-text">Filter</span>
                            </ResponsiveButton>

                            <MenuButton
                                text="Configure"
                                slotProps={{
                                    buttonProps: {
                                        startIcon: (
                                            <SvgIcon style={{ fontSize: 18 }} >
                                                <WrenchScrewdriverIcon />
                                            </SvgIcon>)
                                    }
                                }}>

                                <MenuItem>
                                    <ListItemText>Show Id Column</ListItemText>
                                    <MenuItemToggle value={props.showId} valueChange={(value) => { props.onShowIdChange(value) }} />
                                </MenuItem>

                                <MenuItem>
                                    <ListItemText>Rows per page</ListItemText>
                                    <MenuItemSelect<number>

                                        value={props.contentManagerSearch.state.rowsPerPage}
                                        valueChange={(value) => { props.contentManagerSearch.handleRowsPerPageChange(value) }}
                                        options={[
                                            { label: '10', value: 10 },
                                            { label: '50', value: 50 },
                                            { label: '100', value: 100 }]}
                                    />
                                </MenuItem>

                            </MenuButton>

                        </Stack>

                    </Stack>

                    {!!Object.keys(props.contentManagerSearch.state.filter).length && (<FilterEditor filter={props.contentManagerSearch.state.filter} onChange={props.contentManagerSearch.handleFiltersApply} headCells={headCells} />)}
                </Stack>
            </div>
            {/* <FilterDialog
                filters={filters}
                onApply={handleFiltersApply}
                onClear={handleFiltersClear}
                onClose={filterDialog.dialog.handleClose}
                open={filterDialog.dialog.open}
                properties={filterDialog.filterProperties}
            /> */}
        </>
    );
};

ContentManagerSearch.propTypes = {
    disabled: PropTypes.bool,
    filters: PropTypes.array,
    onFiltersApply: PropTypes.func,
    onFiltersClear: PropTypes.func,
    onQueryChange: PropTypes.func,
    query: PropTypes.string,
    selected: PropTypes.array
};
