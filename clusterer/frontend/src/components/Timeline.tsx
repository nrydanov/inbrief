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
    <TimelineSeparator sx={{ alignItems: 'center', justifyContent: 'center' }} />
    <TimelineContent>
      <Typography 
        variant="subtitle2" 
        sx={{ 
          fontWeight: 600,
          color: 'var(--primary)',
          textTransform: 'uppercase',
          letterSpacing: '0.5px'
        }}
      >
        {date}
      </Typography>
    </TimelineContent>
  </TimelineItem>
);

const EventItem: React.FC<{ time: string; text: string; isLast: boolean }> = ({ time, text, isLast }) => (
  <TimelineItem>
    <TimelineSeparator sx={{ alignItems: 'center', justifyContent: 'center' }}>
      <TimelineDot 
        color="primary" 
        sx={{ 
          width: 16, 
          height: 16,
          boxShadow: '0 0 0 4px var(--selected)'
        }} 
      />
      {!isLast && <TimelineConnector sx={{ bgcolor: 'var(--primary)', width: 2 }} />}
    </TimelineSeparator>
    <TimelineContent>
      <Paper 
        elevation={0}
        sx={{ 
          p: 2,
          mb: 2,
          backgroundColor: 'var(--bg)',
          border: '1px solid',
          borderColor: 'var(--primary-light)',
          borderRadius: 2,
          transition: 'all 0.3s ease',
          '&:hover': {
            backgroundColor: 'var(--hover)',
            transform: 'translateX(4px)'
          }
        }}
      >
        <Typography 
          variant="caption" 
          sx={{ 
            display: 'block',
            color: 'var(--primary)',
            fontWeight: 600,
            mb: 0.5
          }}
        >
          {time}
        </Typography>
        <Typography 
          variant="body2"
          sx={{
            color: 'var(--text-main)',
            lineHeight: 1.6,
            wordBreak: 'break-word'
          }}
        >
          {text}
        </Typography>
      </Paper>
    </TimelineContent>
  </TimelineItem>
);

const Timeline: React.FC<TimelineProps> = ({ events }) => (
  <TimelineMui 
    sx={{ 
      px: 0, 
      py: 0,
      '& .MuiTimelineItem-root': {
        minHeight: 0,
        maxWidth: 600,
        '&:before': {
          display: 'none'
        }
      },
      '& .MuiTimelineSeparator-root': {
        marginLeft: 0,
      },
      '& .MuiTimelineDot-root': {
        marginLeft: 0,
      },
      '& .MuiTimelineContent-root': {
        paddingLeft: 2,
      },
      '& .MuiTimelineConnector-root': {
        marginLeft: 0,
      },
      marginLeft: 0,
      paddingLeft: 0,
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
        <Fade in timeout={500} style={{ transitionDelay: `${idx * 100}ms` }} key={event.datetime + idx}>
          <Box>
            {showDate && <DateLabel date={currentDate} />}
            <EventItem time={time} text={event.text} isLast={idx === events.length - 1} />
          </Box>
        </Fade>
      );
    })}
  </TimelineMui>
);

export default Timeline;
