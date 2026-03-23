**Forward Deployment Engineer**  
**Test**

# **Overview**

This assessment tests your ability to independently integrate a live third-party API into a working application. You will integrate three of Pastelʼs Sigma AI product capabilities: Transaction Monitoring, PEP/Sanctions Screening, and Adverse Media Screening.

The goal is not just to make API calls. You are expected to build a clear, wellstructured application that demonstrates your ability to read technical documentation, translate API capabilities into product value, and present data in a way a non-technical customer could understand. This mirrors the day-today work of a Forward Deployment Engineer at Pastel.

**Sigma exposes three product areas relevant to this assessment:**

	**Transaction Monitoring** — Analyses transactions in real time using AI to detect fraudulent patterns, flag suspicious activity, and return a risk level and recommended action.

	**PEP & Sanctions Screening** — Screens individuals or entities against global Politically Exposed Persons PEP) lists and international sanctions databases OFAC, UN, EU, etc.) to prevent regulated entities from transacting with prohibited parties.

	**Adverse Media Screening** — Searches news and media sources for negative coverage linked to an individual or entity, such as fraud allegations, criminal investigations, or financial misconduct.

**What we are evaluating**

| API comprehension       | Can you independently read documentation and correctly implement an unfamiliar API?                                  |
| :---------------------- | :------------------------------------------------------------------------------------------------------------------- |
| **Data presentation**   | Do you present complex API responses in a way that is clean, structured, and accessible to a non-technical customer? |
| **Attention to detail** | Are all required fields present and correctly formatted? Are edge cases (missing fields, error responses) handled?   |
| **Product**             | Do you go beyond the bare minimum? Does your output help a                                                           |
| **thinking**            | customer understand what the API is telling them and why it matters?                                                 |
| **Communication**       | Is your README and submission clear? Could a Pastel customer success team member understand your demo?               |

# **Your Task**

Build a mock application in any programming language or framework of your choice. The application must integrate Sigmaʼs API and demonstrate all three capabilities described below. Your application should have a frontend interface to interact with.

The Sigma documentation is available at https://sigma-docs.pastel.africa.

You are to make use of the following API credentials:

| API KEY    | 59d01b8c-7d29 4e25 8d93-c5c9f2bb0fdf |
| :--------- | :----------------------------------- |
| API SECRET | ad8a3c3f-1869 4947 94ad-bb216b792824 |

# **Task 1: Transaction Monitoring**

**Your application must:**

Send a transaction to the Sigma Transaction Monitoring API using the instant endpoint

Include realistic, populated data in the a transactionData, device, and anonymizedUserData objects.

Display the following clearly in your applicationʼs output:

1. The transaction details you sent to Sigma (reference, amount, currency, sender, receiver, channel, type, date)

2. The full Sigma response, clearly formatted

3. The final action Sigma recommends (allow / flag / block) — prominently displayed

4. The risk level returned

# **Task 2: PEP & Sanctions Screening**

**Your application must:**

Make a standalone API call to the Sigma PEP and Sanctions endpoint independently.

Pass the details of a information of  individual as required by the documentation

Display the following clearly: a The information you sent to Sigma about the individual  
b The full Sigma response, clearly formatted with the match confidence

# **Task 3: Adverse Media Screening**

**Your application must:**

Make a standalone API call to the Sigma Adverse Media endpoint.

Pass the details of a named individual (can be the same person as Task 2, or a different one).

Display the following clearly: a The information you sent to Sigma about the individual

b The full Sigma response, clearly formatted c A clear summary: were any adverse media findings returned? What category of risk? What sources?

**NOTE You are not required to persist or save the data, you are only required to display the data for viewing purposes only.**

# **Submission**

Deploy your solution online and send in the frontend URL to access it to engineering@pastel.africa on or before Tuesday 24th March, 2026 9 00AM

WAT