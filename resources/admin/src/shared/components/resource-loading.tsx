"use client";

import type { FC } from 'react';
import PropTypes from 'prop-types';
import EllipsisHorizontalIcon from '@heroicons/react/24/outline/EllipsisHorizontalIcon';
import { SxProps } from '@mui/system';
import { SvgIcon, Typography } from '@mui/material';
import { styled } from '@mui/material/styles';

const ResourceLoadingRoot = styled('div')(
  ({ theme }) => ({
    alignItems: 'center',
    // backgroundColor: theme.palette.mode === 'dark'
    //   ? theme.palette.neutral[900]
    //   : theme.palette.neutral[50],
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    padding: theme.spacing(3)
  })
);

interface ResourceLoadingProps {
  message?: string;
  sx?: SxProps;
}

export const ResourceLoading: FC<ResourceLoadingProps> = (props) => {
  const { message, sx } = props;

  return (
    <ResourceLoadingRoot sx={sx}>
      <SvgIcon fontSize="small">
        <EllipsisHorizontalIcon />
      </SvgIcon>
      {message && (
        <Typography
          color="text.secondary"
          sx={{ mt: 2 }}
          variant="body2"
        >
          {message}
        </Typography>
      )}
    </ResourceLoadingRoot>
  );
};

ResourceLoading.propTypes = {
  message: PropTypes.string,
  // @ts-ignore
  sx: PropTypes.object
};
