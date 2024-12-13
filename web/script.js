document.addEventListener("DOMContentLoaded", () => {
    const apiUrl = "http://localhost:8080/api/transactions"; // URL API бэкенда
    const registerUrl = "http://localhost:8080/api/register";
    const loginUrl = "http://localhost:8080/api/login";

    const transactionsList = document.getElementById("transactionList");
    if (!transactionsList) {
        console.error("Не удалось найти элемент для отображения транзакций");
        return;
    }

    const addTransactionForm = document.getElementById("addTransactionForm");
    const registerForm = document.getElementById("registerUserForm");
    const loginForm = document.getElementById("loginUserForm");

    let token = localStorage.getItem('jwtToken'); // Токен JWT из локального хранилища

    // Функция для отображения транзакций
    const displayTransactions = async () => {
        const response = await fetch(apiUrl, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const transactions = await response.json();
        
        // Очищаем текущий список транзакций
        transactionsList.innerHTML = "";
        
        transactions.forEach(transaction => {
            const li = document.createElement("li");
            li.textContent = `${transaction.description} - ${transaction.amount} ${transaction.category}`;
            transactionsList.appendChild(li);
        });
    };

    // Функция для добавления транзакции
    const addTransaction = async (event) => {
        event.preventDefault();
        
        const description = document.getElementById("description").value;
        const amount = parseFloat(document.getElementById("amount").value);
        const category = document.getElementById("category").value;
    
        const transaction = { description, amount, category };
    
        const response = await fetch(apiUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}` // Добавляем токен в заголовок
            },
            body: JSON.stringify(transaction)
        });
    
        if (response.ok) {
            // После успешного добавления, обновляем список транзакций
            displayTransactions();
            addTransactionForm.reset(); // Очищаем форму
        } else {
            alert("Ошибка при добавлении транзакции");
        }
    };

    // Функция для регистрации пользователя
    const registerUser = async (event) => {
        event.preventDefault();

        const email = document.getElementById("registerEmail").value;
        const password = document.getElementById("registerPassword").value;

        const response = await fetch(registerUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            alert("Регистрация прошла успешно!");
            registerForm.reset();
        } else {
            alert("Ошибка при регистрации");
        }
    };

    // Функция для авторизации пользователя
    const loginUser = async (event) => {
        event.preventDefault();

        const email = document.getElementById("loginEmail").value;
        const password = document.getElementById("loginPassword").value;

        const response = await fetch(loginUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();
            token = data.token; // Сохраняем токен
            localStorage.setItem('jwtToken', token); // Сохраняем токен в локальном хранилище
            alert("Авторизация прошла успешно!");
            displayTransactions();
        } else {
            alert("Ошибка при авторизации");
        }
    };

    // Загружаем транзакции при старте, если пользователь авторизован
    if (token) {
        displayTransactions();
    }

    // Обработчики форм
    registerForm.addEventListener("submit", registerUser);
    loginForm.addEventListener("submit", loginUser);
    addTransactionForm.addEventListener("submit", addTransaction);
});
