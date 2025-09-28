# Builder Maven desde GHCR (variante temurin-8)
FROM ghcr.io/library/maven:3.8.7-eclipse-temurin-8 AS build
WORKDIR /src
COPY pom.xml .
COPY src ./src
RUN mvn -B -DskipTests package

# Runtime Temurin JRE 8 desde GHCR (evita Docker Hub)
FROM ghcr.io/adoptium/temurin:8-jre
WORKDIR /app
COPY --from=build /src/target/*jar /app/app.jar

# Tu servicio escucha 8083 (en compose pasas SERVER_PORT=8083)
EXPOSE 8083
ENTRYPOINT ["java","-jar","/app/app.jar"]
