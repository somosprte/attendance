import React, { useState } from 'react';
import axios from 'axios';

const App = () => {
  const [meeting, setMeeting] = useState({ title: '', date: '', time: '', description: '' });
  const [meetingID, setMeetingID] = useState('');
  const [attendees, setAttendees] = useState([]);
  const [name, setName] = useState('');

  const createMeeting = () => {
    axios.post('/api/meetings', meeting)
      .then((response) => {
        setMeetingID(response.data.id);
        alert('Meeting created! ID: ' + response.data.id);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const registerAttendee = () => {
    axios.post(`/api/register/${meetingID}?name=${name}`)
      .then((response) => {
        setAttendees(response.data.attendees);
        alert('Successfully registered!');
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const getMeeting = () => {
    axios.get(`/api/meetings/${meetingID}`)
      .then((response) => {
        setAttendees(response.data.attendees);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return (
    <div>
      <h1>Create Meeting</h1>
      <input type="text" placeholder="Title" onChange={(e) => setMeeting({ ...meeting, title: e.target.value })} />
      <input type="date" onChange={(e) => setMeeting({ ...meeting, date: e.target.value })} />
      <input type="time" onChange={(e) => setMeeting({ ...meeting, time: e.target.value })} />
      <input type="text" placeholder="Description" onChange={(e) => setMeeting({ ...meeting, description: e.target.value })} />
      <button onClick={createMeeting}>Create Meeting</button>
      <h1>Register for Meeting</h1>
      <input type="text" placeholder="Meeting ID" onChange={(e) => setMeetingID(e.target.value)} />
      <input type="text" placeholder="Your Name" onChange={(e) => setName(e.target.value)} />
      <button onClick={registerAttendee}>Register</button>
      <button onClick={getMeeting}>View Attendees</button>
      <h1>Attendees</h1>
      <ul>
        {attendees.map((attendee, index) => <li key={index}>{attendee}</li>)}
      </ul>
    </div>
  );
};

export default App;
