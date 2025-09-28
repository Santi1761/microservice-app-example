package com.elgris.usersapi.security;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.GenericFilterBean;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

@Component
public class JwtAuthenticationFilter extends GenericFilterBean {

    @Value("${jwt.secret}")
    private String jwtSecret;

    @Override
    public void doFilter(final ServletRequest req, final ServletResponse res, final FilterChain chain)
            throws IOException, ServletException {

        final HttpServletRequest request = (HttpServletRequest) req;
        final HttpServletResponse response = (HttpServletResponse) res;

        if ("OPTIONS".equalsIgnoreCase(request.getMethod())) {
            response.setStatus(HttpServletResponse.SC_OK);
            chain.doFilter(req, res);
            return;
        }

        final String uri = request.getRequestURI();
        final boolean isPublicUsersGet = "GET".equalsIgnoreCase(request.getMethod())
                && uri != null
                && (uri.equals("/users") || uri.equals("/users/") || uri.startsWith("/users/"));

        final String authHeader = request.getHeader("authorization");

        if (isPublicUsersGet && (authHeader == null || !authHeader.startsWith("Bearer "))) {
            chain.doFilter(req, res);
            return;
        }

        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Missing or invalid Authorization header");
            return;
        }

        final String token = authHeader.substring(7);

        try {
            final Claims claims = Jwts.parser()
                    .setSigningKey(jwtSecret.getBytes())
                    .parseClaimsJws(token)
                    .getBody();

            request.setAttribute("claims", claims);
            chain.doFilter(req, res);

        } catch (final SignatureException e) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Invalid token signature");
        } catch (final Exception e) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Token validation failed: " + e.getMessage());
        }
    }
}
