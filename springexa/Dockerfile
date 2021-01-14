FROM openjdk:11

EXPOSE 8080

ARG JAR_FILE=target/springexa.jar
ADD ${JAR_FILE} springexa.jar

ENTRYPOINT ["java", "-jar", "/springexa.jar"]