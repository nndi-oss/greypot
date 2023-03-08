import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  let testData = {
    Name: "k6-test.html",
    Template: `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Sample Report</title>
</head>
<body>
    <h1>Sample User Report</h1>
    <div>Generated at {{ data.generatedAt }} </div>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Age</th>
                <th>Email</th>
            </tr>
        </thead>
        <tbody>
            {% for u in data.users %}
            <tr>
                <td>{{ u.name }}</td>
                <td>{{ u.age }}</td>
                <td>{{ u.email }}</td>
            </tr>
            {% endfor %}
        </tbody>
    </table>
</body>
</html>
`,
    Data: {
        generatedAt: (new Date()).toISOString(),
        users: [
            { name: "John Doe", age: 30, email: "john@example.com" },
            { name: "Mary Doe", age: 66, email: "mary@example.com" },
            { name: "Bob Doe", age: 41, email: "bob@example.com" },
            { name: "Jane Doe", age: 26, email: "jane@example.com" }
        ]
    }
  }

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post('http://localhost:7665/_studio/generate/pdf/k6test', JSON.stringify(testData), params);
  
  sleep(1)
}
