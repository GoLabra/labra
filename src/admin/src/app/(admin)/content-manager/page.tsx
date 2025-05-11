"use client";

import { useDocumentTitle } from "@/hooks/use-document-title";
import {
  Container,
  Stack,
  Box,
  Card,
  CardContent,
  Typography,
  styled,
} from "@mui/material";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import { useCallback, useEffect, useState } from "react";
import { OfficialCommunicationChannels } from "@/shared/components/official-communication-channels";

const BulletItem = styled(Stack)(({ theme }) => ({
  border: "2px dashed #80808024",
  borderRadius: "3px",
  padding: "10px",
  zIndex: 1,
  backgroundColor: "var(--mui-palette-background-default)",
}));

export default function ContentManager() {
  useDocumentTitle({ title: "Content Manager" });

  const [rotation, setRotation] = useState(0);

  const setRandomRotation = useCallback(() => {
    setRotation(Math.random() * 5 - 2);
  }, [setRotation]);
  
  useEffect(() => {
    setRandomRotation();
  }, [setRandomRotation]);

  return (
    <Box
      sx={{
        flexGrow: 1,
        py: 4,
      }}
    >
      <Container maxWidth="md">
        <Box
          sx={{
            flexGrow: 1,
            py: 4,
          }}
        >
          <Stack spacing={5}>
            <Typography color="text.primary" variant="body2">
              Manage your content effortlessly by selecting an entity from the
              sidebar. This page is designed to help you oversee, organize, and
              maintain all your entries effectively.
            </Typography>

            <Stack spacing={2}>
              <BulletItem alignItems="center" direction="row" spacing={3}>
                <Typography>
                  <FiberManualRecordIcon />
                </Typography>

                <div>
                  <Typography variant="subtitle1" color="text.primary">
                    Content Overview
                  </Typography>
                  <Typography color="text.secondary" variant="body2">
                    View a comprehensive list of entries related to your
                    selected entity, with important details such as titles,
                    modification dates, and statuses.
                  </Typography>
                </div>
              </BulletItem>

              <BulletItem
                alignItems="center"
                direction="row"
                spacing={3}
                sx={{
                  zIndex: 2,
                  transform: `rotate(${rotation}deg)`,
                  transition: "transform 0.2s ease-in-out",
                }}
                onClick={setRandomRotation}
              >
                <Typography>
                  <FiberManualRecordIcon />
                </Typography>

                <div>
                  <Typography variant="subtitle1" color="text.primary">
                    Flexible Management
                  </Typography>

                  <Typography color="text.secondary" variant="body2">
                    Easily update or remove content to keep everything accurate
                    and relevant.
                  </Typography>
                </div>
              </BulletItem>

              <BulletItem alignItems="center" direction="row" spacing={3}>
                <Typography>
                  <FiberManualRecordIcon />
                </Typography>

                <div>
                  <Typography variant="subtitle1" color="text.primary">
                    Advanced Tools
                  </Typography>

                  <Typography color="text.secondary" variant="body2">
                    Utilize filtering, sorting, and ordering options to find and
                    manage entries efficiently.
                  </Typography>
                </div>
              </BulletItem>

              <BulletItem alignItems="center" direction="row" spacing={3}>
                <Typography>
                  <FiberManualRecordIcon />
                </Typography>

                <div>
                  <Typography variant="subtitle1" color="text.primary">
                    Streamlined Interaction
                  </Typography>

                  <Typography color="text.secondary" variant="body2">
                    Perform content actions quickly and seamlessly for an
                    optimized workflow.
                  </Typography>
                </div>
              </BulletItem>
            </Stack>

            <OfficialCommunicationChannels></OfficialCommunicationChannels>
          </Stack>
        </Box>
      </Container>
    </Box>
  );
}
