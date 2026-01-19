import { useEffect, useState } from 'react'
import './App.css'
import type { User } from './types/user';

function App() {
  const [currentuser, setCurrentUser] = useState<TelegramUser>();
  const [users, setUsers] = useState<User[]>();

  useEffect(() => {
    const tg = window.Telegram?.WebApp;
    if (tg && tg.initDataUnsafe?.user) {
      setTimeout(() => setCurrentUser(tg.initDataUnsafe.user), 0);
      console.log(currentuser);
    } else {
      console.warn("⚠️ Telegram WebApp API недоступен. Запусти приложение через Telegram.");
    }
  }, []);

  const handleGetUsers = () => {
    fetch('http://localhost:3000/api/users')
      .then(response => response.json())
      .then(data => {
        setUsers(data)
      })
      
      .catch(error => console.error('Ошибка при получении данных о пользователе:', error));
  };

  return (
    <div>
    <button onClick={handleGetUsers}>получить пользователей</button>
      <h1>FitBoard Mini App</h1>
      {users && users.length > 0 ? (
        <div>
          {users.map(user => (
            <div key={user.telegram_id}>
              <p>ID: {user.telegram_id}</p>
              <p>role: {user.role}</p>
              <p>Username: {user.username}</p>
              <p>First name: {user.first_name}</p>
              <p>Last name: {user.last_name}</p>
            </div>
          ))}
        </div>
      ) : (
        <p>Нет данных о пользователе</p>
      )}
    </div>
  );
}

export default App
