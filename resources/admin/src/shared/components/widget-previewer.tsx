"use client";

import type { FC, ReactNode } from 'react';
import PropTypes from 'prop-types';
import type { SxProps } from '@mui/system';
import { Card, Stack, Typography } from '@mui/material';

interface WidgetPreviewerProps {
  children: ReactNode;
  description?: string | ReactNode;
  title?: string | ReactNode;
  sx?: SxProps;
}

export const WidgetPreviewer: FC<WidgetPreviewerProps> = (props) => {
  const { title, description, children, ...other } = props;

  return (
    <Stack
      spacing={2}
      {...other}
    >
      <Stack spacing={1}>
        {
          typeof title === 'string'
            ? (
              <Typography variant="subtitle1">
                {title}
              </Typography>
            )
            : title
        }
        {
          typeof description === 'string'
            ? (
              <Typography
                color="text.secondary"
                variant="body2"
              >
                {description}
              </Typography>
            )
            : description
        }
      </Stack>
      <Card>
        {children}
      </Card>
    </Stack>
  );
};