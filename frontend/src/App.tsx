import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [user, setUser] = useState<TelegramUser>();

  useEffect(() => {
    const tg = window.Telegram?.WebApp;
    if (tg && tg.initDataUnsafe?.user) {
      setTimeout(() => setUser(tg.initDataUnsafe.user), 0);
    } else {
      console.warn("⚠️ Telegram WebApp API недоступен. Запусти приложение через Telegram.");
    }
  }, []);

  return (
    <div>
      <h1>Telegram Mini App</h1>
      {user ? (
        <div>
          <p>ID: {user.id}</p>
          <p>Username: {user.username}</p>
          <p>First name: {user.first_name}</p>
          <p>Last name: {user.last_name}</p>
        </div>
      ) : (
        <p>Нет данных о пользователе</p>
      )}
    </div>
  );
}

export default App
