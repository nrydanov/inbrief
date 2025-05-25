import React from 'react';
import TimelineMui from '@mui/lab/Timeline';
import TimelineItem from '@mui/lab/TimelineItem';
import TimelineSeparator from '@mui/lab/TimelineSeparator';
import TimelineConnector from '@mui/lab/TimelineConnector';
import TimelineContent from '@mui/lab/TimelineContent';
import TimelineDot from '@mui/lab/TimelineDot';
import { Typography, Paper, Box, Fade } from '@mui/material';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export interface TimelineEvent {
  datetime: string;
  text: string;
}

interface TimelineProps {
  events: TimelineEvent[];
}

const DateLabel: React.FC<{ date: string }> = ({ date }) => (
  <TimelineItem>
    <TimelineSeparator />
    <TimelineContent>
      <Typography
        variant="subtitle2"
        sx={{
          fontWeight: 600,
          color: 'primary.main',
          textTransform: 'uppercase',
          letterSpacing: '0.5px',
        }}
      >
        {date}
      </Typography>
    </TimelineContent>
  </TimelineItem>
);

const EventItem: React.FC<{ time: string; text: string; isLast: boolean }> = ({
  time,
  text,
  isLast,
}) => (
  <TimelineItem>
    <TimelineSeparator>
      <TimelineDot
        color="primary"
        sx={{
          width: 16,
          height: 16,
          boxShadow: '0 0 0 4px rgba(33, 150, 243, 0.1)',
        }}
      />
      {!isLast && <TimelineConnector sx={{ bgcolor: 'primary.main', width: 2 }} />}
    </TimelineSeparator>
    <TimelineContent>
      <Paper
        elevation={0}
        sx={{
          p: 2,
          mb: 2,
          backgroundColor: 'rgba(33, 150, 243, 0.04)',
          border: '1px solid',
          borderColor: 'primary.light',
          borderRadius: 2,
          transition: 'all 0.3s ease',
          '&:hover': {
            backgroundColor: 'rgba(33, 150, 243, 0.08)',
            transform: 'translateX(4px)',
          },
        }}
      >
        <Typography
          variant="caption"
          sx={{
            display: 'block',
            color: 'primary.main',
            fontWeight: 600,
            mb: 0.5,
          }}
        >
          {time}
        </Typography>
        <Typography
          variant="body2"
          sx={{
            color: 'text.primary',
            lineHeight: 1.6,
            wordBreak: 'break-word',
          }}
        >
          {text}
        </Typography>
      </Paper>
    </TimelineContent>
  </TimelineItem>
);

const Timeline: React.FC<TimelineProps> = ({ events }) => {
  if (!events || events.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" align="center">
        Нет событий для отображения
      </Typography>
    );
  }

  return (
    <TimelineMui
      sx={{
        px: 0,
        py: 0,
        marginLeft: 0,
        paddingLeft: 0,
        '& .MuiTimelineItem-root': {
          minHeight: 0,
          maxWidth: 600,
          '&:before': { display: 'none' },
        },
        '& .MuiTimelineContent-root': { paddingLeft: 2 },
      }}
    >
      {events.map((event, idx) => {
        const dateObj = new Date(event.datetime);
        const currentDate = format(dateObj, 'd MMMM yyyy', { locale: ru });
        const time = format(dateObj, 'HH:mm', { locale: ru });
        const prevEvent = events[idx - 1];
        const prevDate = prevEvent ? format(new Date(prevEvent.datetime), 'd MMMM yyyy', { locale: ru }) : null;
        const showDate = currentDate !== prevDate;

        return (
          <Fade in timeout={500} style={{ transitionDelay: `${idx * 100}ms` }} key={`${event.datetime}-${idx}`}>
            <Box>
              {showDate && <DateLabel date={currentDate} />}
              <EventItem time={time} text={event.text} isLast={idx === events.length - 1} />
            </Box>
          </Fade>
        );
      })}
    </TimelineMui>
  );
};

export default Timeline;
