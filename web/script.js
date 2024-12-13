document.addEventListener("DOMContentLoaded", () => {
    const apiUrl = "http://localhost:8080/api/transactions"; // URL API бэкенда

    const transactionsList = document.getElementById("transactionList");
    if (!transactionsList) {
        console.error("Не удалось найти элемент для отображения транзакций");
        return;
    }

    const addTransactionForm = document.getElementById("addTransactionForm");

    // Функция для отображения транзакций
    const displayTransactions = async () => {
        const response = await fetch(apiUrl);
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
                "Content-Type": "application/json"
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

    // Загружаем транзакции при старте
    displayTransactions();

    // Обработчик формы
    addTransactionForm.addEventListener("submit", addTransaction);
});
