// "use client";

// import type { FC } from 'react';
// import PropTypes from 'prop-types';
// import type { DropzoneOptions } from 'react-dropzone';
// import { useDropzone } from 'react-dropzone';
// import ArrowUpOnSquareIcon from '@heroicons/react/24/outline/ArrowUpOnSquareIcon';
// import type { SxProps } from '@mui/system';
// import { Avatar, Box, SvgIcon, Typography } from '@mui/material';

// interface FileDropzoneProps extends DropzoneOptions {
//   caption?: string;
//   sx?: SxProps;
// }

// export const FileDropzone: FC<FileDropzoneProps> = (props) => {
//   const { accept, caption, maxFiles, maxSize, minSize, onDrop, sx } = props;
//   const { getRootProps, getInputProps, isDragActive } = useDropzone({
//     accept,
//     maxFiles,
//     maxSize,
//     minSize,
//     onDrop
//   });

//   return (
//     <Box
//       sx={{
//         alignItems: 'center',
//         borderColor: (theme) => theme.palette.mode === 'dark'
//           ? 'neutral.800'
//           : 'neutral.200',
//         borderRadius: 1,
//         borderStyle: 'dashed',
//         borderWidth: 1,
//         cursor: 'pointer',
//         display: 'flex',
//         flexDirection: 'column',
//         flexWrap: 'wrap',
//         justifyContent: 'center',
//         outline: 'none',
//         width: '100%',
//         py: 2,
//         ...(isDragActive && {
//           borderColor: 'primary.main',
//           backgroundColor: 'action.hover'
//         }),
//         '&:hover': {
//           borderColor: 'primary.main',
//           backgroundColor: 'action.hover'
//         },
//         ...sx
//       }}
//       {...getRootProps()}
//     >
//       <Avatar
//         sx={{
//           backgroundColor: 'neutral.200',
//           color: 'text.secondary',
//           height: 42,
//           width: 42
//         }}
//       >
//         <SvgIcon fontSize="small">
//           <ArrowUpOnSquareIcon />
//         </SvgIcon>
//       </Avatar>
//       {caption && (
//         <Typography
//           align="center"
//           color="text.secondary"
//           sx={{ mt: 1 }}
//           variant="caption"
//         >
//           {caption}
//         </Typography>
//       )}
//       <input {...getInputProps()} />
//     </Box>
//   );
// };

// FileDropzone.propTypes = {
//   // @ts-ignore
//   accept: PropTypes.objectOf(PropTypes.arrayOf(PropTypes.string)),
//   caption: PropTypes.string,
//   maxFiles: PropTypes.number,
//   maxSize: PropTypes.number,
//   minSize: PropTypes.number,
//   onDrop: PropTypes.func,
//   sx: PropTypes.object
// };
