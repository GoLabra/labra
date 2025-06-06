import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import { fileData } from '../file-thumbnail/utils';
import { FileRejection } from 'react-dropzone';
import { fData } from '@/lib/utils/format-number';

// import { fData } from 'src/utils/format-number';

// import { varAlpha } from 'src/theme/styles';

// import { fileData } from '../../file-thumbnail';

// ----------------------------------------------------------------------

interface RejectionFile {
    files: readonly FileRejection[];
}
export function RejectionFiles({ files }: RejectionFile) {
  if (!files.length) {
    return null;
  }

  return (
    <Paper
      variant="outlined"
      sx={{
        py: 1,
        px: 2,
        mt: 3,
        textAlign: 'left',
        borderStyle: 'dashed',
        borderColor: 'error.main',
        //bgcolor: (theme) => varAlpha(theme.vars.palette.error.mainChannel, 0.08),
      }}
    >
      {files.map(({ file, errors }) => {
        const { path, size } = fileData(file);

        return (
          <Box key={path} sx={{ my: 1 }}>
            <Typography variant="subtitle2" noWrap>
              {path} - {size ? fData(size) : ''}
            </Typography>

            {errors.map((error) => (
              <Box key={error.code} component="span" sx={{ typography: 'caption' }}>
                - {error.message}
              </Box>
            ))}
          </Box>
        );
      })}
    </Paper>
  );
}
