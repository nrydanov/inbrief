import React from 'react';
import TimelineMui from '@mui/lab/Timeline';
import TimelineItem from '@mui/lab/TimelineItem';
import TimelineSeparator from '@mui/lab/TimelineSeparator';
import TimelineConnector from '@mui/lab/TimelineConnector';
import TimelineContent from '@mui/lab/TimelineContent';
import TimelineDot from '@mui/lab/TimelineDot';
import { Typography } from '@mui/material';

export interface TimelineEvent {
  datetime: string;
  text: string;
}

interface TimelineProps {
  events: TimelineEvent[];
}

<TimelineMui
  sx={{
    px: 0,
    py: 0,
    '& .MuiTimelineItem-root': {
      minHeight: 0,
    },
    '& .MuiTimelineSeparator-root': {
      marginLeft: 0,
    },
    '& .MuiTimelineDot-root': {
      marginLeft: 0,
    },
    '& .MuiTimelineContent-root': {
      paddingLeft: 2, // Optional: add space between dot and content
    },
    '& .MuiTimelineConnector-root': {
      marginLeft: 0,
    },
    // Remove padding from the timeline itself
    marginLeft: 0,
    paddingLeft: 0,
  }}
>
  {/* timeline items */}
</TimelineMui>


const DateLabel: React.FC<{ date: string }> = ({ date }) => (
  <TimelineItem>
    <TimelineSeparator>
      <TimelineDot sx={{ width: 6, height: 6, bgcolor: 'primary.main' }} />
      <TimelineConnector sx={{ bgcolor: 'primary.main', width: 2 }} />
    </TimelineSeparator>
    <TimelineContent>
      <Typography variant="caption" sx={{ fontWeight: 'bold' }}>{date}</Typography>
    </TimelineContent>
  </TimelineItem>
);

const EventItem: React.FC<{ time: string; text: string; isLast: boolean }> = ({ time, text, isLast }) => (
  <TimelineItem>
    <TimelineSeparator>
      <TimelineDot color="primary" sx={{ width: 12, height: 12 }} />
      {!isLast && <TimelineConnector sx={{ bgcolor: 'primary.main', width: 2 }} />}
    </TimelineSeparator>
    <TimelineContent>
      <Typography variant="body2" color="text.secondary">{time}</Typography>
      <Typography variant="body2">{text}</Typography>
    </TimelineContent>
  </TimelineItem>
);

const Timeline: React.FC<TimelineProps> = ({ events }) => {
  return (
    <TimelineMui sx={{ px: 0, py: 0 }}>
      {events.map((event, idx) => {
        const dateObj = new Date(event.datetime);
        const currentDate = dateObj.toLocaleDateString('ru-RU');
        const time = dateObj.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
        const prevEvent = events[idx - 1];
        const prevDate = prevEvent ? new Date(prevEvent.datetime).toLocaleDateString('ru-RU') : null;
        const showDate = currentDate !== prevDate;
        return (
          <React.Fragment key={event.datetime + idx}>
            {showDate && <DateLabel date={currentDate} />}
            <EventItem time={time} text={event.text} isLast={idx === events.length - 1} />
          </React.Fragment>
        );
      })}
    </TimelineMui>
  );
};

export default Timeline;
