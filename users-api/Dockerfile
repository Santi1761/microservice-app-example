# Builder: OK con el mirror de Google
FROM mirror.gcr.io/library/maven:3.8.7-eclipse-temurin-8 AS build
WORKDIR /src
COPY pom.xml .
COPY src ./src
RUN mvn -B -DskipTests package

# ✅ Runtime: usa OpenJDK 8 desde el mirror de Google (sin rate-limit ni GHCR)
FROM mirror.gcr.io/library/openjdk:8-jre
WORKDIR /app
COPY --from=build /src/target/*jar /app/app.jar

# Tu compose expone users-api por 8083
EXPOSE 8083
ENTRYPOINT ["java","-jar","/app/app.jar"]
