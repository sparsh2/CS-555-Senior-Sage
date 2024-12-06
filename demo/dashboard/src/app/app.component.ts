import { Component, NO_ERRORS_SCHEMA } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AccessLogsComponent } from './access-logs/access-logs.component';
import { ChatHistoryComponent } from './chat-history/chat-history.component';
import { HeaderComponent } from './header/header.component';
import { RemindersComponent } from './reminders/reminders.component';
import { PreferencesComponent } from './preferences/preferences.component';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import {FormsModule} from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {MatFormFieldModule} from '@angular/material/form-field';
import { MatSelectChange } from '@angular/material/select';
import axios from 'axios';

interface User {
  name: string;
}


interface DataRequestBody {
  // Define the structure of your request body here
  user_id: string;
  requester_token: string;
}

interface LogsRequestBody {
  user_id: string;
}

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet, AccessLogsComponent, ChatHistoryComponent, HeaderComponent, RemindersComponent, PreferencesComponent, MatTableModule, CommonModule,
    FormsModule, MatInputModule, MatSelectModule, MatFormFieldModule
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  schemas: [NO_ERRORS_SCHEMA]
})
export class AppComponent {
  title = 'dashboard';
  currentUser = 'Alice';

  displayUsers: User[] = [
    {name: 'Alice'},
    {name: 'Bob'}
  ];

  usersList: Map<string, string[]> = new Map();
  displayedColumns: string[] = ['reminder_for', 'time', 'frequency', 'cron_job'];
  
  data: any = {}
  logs: any = []

  async getData() {
    try {
      let req2: LogsRequestBody = {
        user_id: this.usersList.get(this.currentUser)![0]
      }
      const response2 = await axios.post('http://127.0.0.1:52714/access-logs', req2);
      this.logs = response2.data.logs;

      let req: DataRequestBody = {
        user_id: this.usersList.get(this.currentUser)![0],
        requester_token: this.usersList.get(this.currentUser)![1]
      }
      const response = await axios.post('http://127.0.0.1:52714/data', req);
      console.log('Data fetched:', response.data);
      this.data = response.data;

      
    } catch (error) {
      console.error('Error fetching data:', error);
      // return null;
    }
  }

  ngOnInit() {
    this.usersList.set('Alice', ['91c6e005d9fcd20a6cbf1fe1dfdd319d', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzYWdlLXNlcnZlciIsImV4cCI6MTc2OTQ1MTQ4MSwibmJmIjoxNzMzNDU1MDgxLCJpYXQiOjE3MzM0NTUwODEsInVzZXJfaWQiOiJmYWY5YWUzOGQ4NmFlMDcxOWFkODFhZmMxOThkYjFiNyJ9.i8X_bzNg8--r8BuRfU7xGqHmczhrfY9RX0AOE3ojvkw'])
    this.usersList.set('Bob', ['ee82c18d386dd4f0b6e7fae47e5fc1a6', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzYWdlLXNlcnZlciIsImV4cCI6MTc2OTQ1MTQ4MSwibmJmIjoxNzMzNDU1MDgxLCJpYXQiOjE3MzM0NTUwODEsInVzZXJfaWQiOiJmYWY5YWUzOGQ4NmFlMDcxOWFkODFhZmMxOThkYjFiNyJ9.i8X_bzNg8--r8BuRfU7xGqHmczhrfY9RX0AOE3ojvkw'])
    this.getData();
  }

}
/*
{
    "reminder_details": [],
    "rpm_readings": [],
    "preferences": [],
    "chat_history": [],
    "question_responses": [],
    "voice_selection": "Nova",
    "name": "Alice",
    "msg": "data fetched successfully",
    "question_counts": {}
}
*/