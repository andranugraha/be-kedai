# Kedai Backend

## A. Profiling & Benchmarking

1. Load & performance test
    - These tests are used to evaluate the API's ability to handle varying levels of traffic and ensure that it can perform well under high load conditions. 
      - Burst traffic tests involve sending a sudden spike of traffic to the API to simulate a sudden surge in usage.
      - High traffic tests involve sustained high levels of traffic over an extended period to evaluate the API's ability to handle continuous usage.
      - Sustained traffic tests involve sending a steady stream of traffic over an extended period to evaluate the API's long-term performance.
      - Overall, these tests help identify any bottlenecks or performance issues that may arise under different usage scenarios and allow developers to make improvements to optimize the API's performance.

2. Profile
    - Run profiling tool to see app's profile
    - Profiling is a software development technique used to measure the performance of an application or system. It involves collecting and analyzing data about the execution of the application or system to identify performance bottlenecks, memory leaks, and other issues that may impact its overall performance.
    - Profiling can be used for a variety of purposes, including:
      - Performance Optimization: Profiling can help developers identify parts of their code that are slowing down the application and optimize them for better performance.
      - Memory Optimization: Profiling can identify memory leaks and excessive memory usage, allowing developers to optimize their code and reduce memory usage.
      - Debugging: Profiling can help identify errors or unexpected behavior in an application, making it easier for developers to locate and fix bugs.
      - Capacity Planning: Profiling can provide insights into how an application or system will perform under different levels of load, helping developers plan for scalability and capacity needs.
      - In summary, profiling is a valuable tool for software developers to improve the performance, reliability, and scalability of their applications and systems.
    - Here are some steps to follow when using profiling to identify these issues:
      1. Identify the problem: Start by identifying the specific performance issue you want to investigate, such as slow response times, high CPU usage, or excessive memory usage.
      2. Choose a profiling tool: Select a profiling tool that is appropriate for your application or system, and that can provide the data you need to diagnose the problem.
      3. Configure the profiling tool: Configure the profiling tool to collect the relevant performance data, such as CPU usage, memory usage, or network activity.
      4. Run the application: Start the application and run it under normal operating conditions, making sure to capture a representative sample of usage data.
      5. Analyze the data: Use the profiling tool to analyze the data collected during the application's execution, looking for patterns or trends that indicate performance bottlenecks or memory leaks.
      6. Identify the root cause: Identify the root cause of the performance issue by examining the data and looking for areas of the code or system that may be causing the problem.
      7. Optimize or fix the code: Once you have identified the root cause of the problem, optimize or fix the code or system to address the issue and improve performance.
      8. By following these steps, you can use profiling tools to diagnose performance bottlenecks, memory leaks, and other issues that may impact your application's overall performance.

3. Benchmark
    - Benchmarking is the process of measuring the performance of a system or application by comparing it to a standard or reference point. In software development, benchmarking is often used to compare the performance of different programming languages, libraries, algorithms, or hardware configurations.
    - The benchmarking process typically involves the following steps:
      1. Define the test criteria: Define the criteria that will be used to measure the performance of the system or application, such as response time, throughput, or latency.
      2. Select the benchmark tool: Select a benchmarking tool that is appropriate for the type of system or application being tested, and that can provide the data needed to evaluate performance.
      3. Run the benchmark: Run the benchmarking tool against the system or application under test, using the test criteria to collect data on its performance.
      4. Analyze the results: Analyze the data collected by the benchmarking tool to evaluate the system or application's performance against the standard or reference point.
      5. Draw conclusions and take action: Use the results of the benchmarking test to draw conclusions about the performance of the system or application, and take action to optimize or improve its performance as needed.
      6. Overall, benchmarking is a valuable technique for software developers to measure and compare the performance of different systems and applications, helping to identify areas for optimization and improvement.
    - Benchmark output: BenchmarkFunctionName/SubtestName-NumberofThreads | NumberofIterations | TimePerOperation