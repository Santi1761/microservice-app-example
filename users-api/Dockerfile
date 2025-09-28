# Builder Maven desde el mirror de Google (variante temurin-8)
FROM mirror.gcr.io/library/maven:3.8.7-eclipse-temurin-8 AS build
WORKDIR /src
COPY pom.xml .
COPY src ./src
RUN mvn -B -DskipTests package

# Runtime Temurin JRE 8 desde Adoptium en GHCR (este sí es público en GHCR)
FROM ghcr.io/adoptium/temurin:8-jre
WORKDIR /app
COPY --from=build /src/target/*jar /app/app.jar

# Tu compose mapea 8083 para users-api
EXPOSE 8083
ENTRYPOINT ["java","-jar","/app/app.jar"]
