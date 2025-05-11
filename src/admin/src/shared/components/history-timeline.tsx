
import Timeline from '@mui/lab/Timeline';
import TimelineItem, { timelineItemClasses } from '@mui/lab/TimelineItem';
import TimelineSeparator from '@mui/lab/TimelineSeparator';
import TimelineConnector from '@mui/lab/TimelineConnector';
import TimelineContent from '@mui/lab/TimelineContent';
import TimelineDot from '@mui/lab/TimelineDot';

import { Stack, styled, Typography } from "@mui/material";
import { evalManifestWithRetries } from "next/dist/server/load-components";
import { ReactNode } from "react";
import TimelineOppositeContent, { timelineOppositeContentClasses } from '@mui/lab/TimelineOppositeContent';

export type HistoryEvent = {
    label: string;
    details: ReactNode;
}

interface HistoryTimeLineProps {
    events?: HistoryEvent[];
}





const TimelineConnectorNoSpace = styled(TimelineConnector)(({ theme }) => ({

    "&:before": {
        content: '" "',
        position: 'absolute',
        width: '2px',
        height: '14px',
        backgroundColor: '#afafaf',
        top: '22px'
    },
    "&:after": {
        content: '" "',
        position: 'absolute',
        width: '2px',
        height: '14px',
        backgroundColor: '#afafaf',
        bottom: '-13px'
    },
}));


export const HistoryTimeLine = (props: HistoryTimeLineProps) => {

    const { events } = props;

    if(!events || !events.length){
        return null;
    }

    return (
        <Timeline
        sx={{
          [`& .${timelineOppositeContentClasses.root}`]: {
            flex: 0.3,
          },
        }}
      >

            {events.map((ev, index) => (
                <TimelineItem 
                    key={ev.label}
                    sx={{
                        ...(index == events.length - 1 && {
                            minHeight: 0
                        })
                    }}>
                           <TimelineOppositeContent color="text.secondary">
                           <Typography
                                color="text.secondary"
                                fontSize={12}
                                lineHeight='1.5rem'>
                                {ev.label}
                            </Typography>
                            </TimelineOppositeContent>

                    <TimelineSeparator>
                        <TimelineDot />
                        {index < events.length-1 && <TimelineConnectorNoSpace />}
                    </TimelineSeparator>
                    <TimelineContent>
                        {ev.details}
                    </TimelineContent>
                </TimelineItem>
            ))}
        </Timeline>
    )
}