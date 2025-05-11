"use client";

import type { FC, ReactNode } from 'react';
import PropTypes from 'prop-types';
import type { ListItemProps } from '@mui/material';
import { Box, ListItem, ListItemText, Typography } from '@mui/material';

interface PropertyListItemProps extends ListItemProps {
  align?: string;
  children?: ReactNode;
  component?: any;
  label: string;
  value?: string;
}

export const PropertyListItem: FC<PropertyListItemProps> = (props) => {
  const {
    align = 'vertical',
    children,
    component,
    label,
    value = '',
    ...other
  } = props;

  return (
    <ListItem
      component={component}
      disableGutters
      sx={{
        px: 3,
        py: 1.5
      }}
      {...other}
    >
      <ListItemText
        disableTypography
        primary={(
          <Typography
            sx={{ minWidth: align === 'vertical' ? 'inherit' : 180 }}
            variant="subtitle2"
          >
            {label}
          </Typography>
        )}
        secondary={(
          <Box
            sx={{
              flexGrow: 1,
              mt: align === 'vertical' ? 0.5 : 0
            }}
          >
            {children || (
              <Typography
                color="text.secondary"
                variant="body2"
              >
                {value}
              </Typography>
            )}
          </Box>
        )}
        sx={{
          alignItems: 'flex-start',
          display: 'flex',
          flexDirection: align === 'vertical' ? 'column' : 'row',
          my: 0
        }}
      />
    </ListItem>
  );
};

PropertyListItem.propTypes = {
  align: PropTypes.string,
  component: PropTypes.any,
  label: PropTypes.string.isRequired,
  value: PropTypes.string
};
