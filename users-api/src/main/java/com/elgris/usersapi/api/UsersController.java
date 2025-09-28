package com.elgris.usersapi.api;

import com.elgris.usersapi.models.User;
import com.elgris.usersapi.repository.UserRepository;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.util.LinkedList;
import java.util.List;

@RestController()
@RequestMapping("/users")
public class UsersController {

    @Autowired
    private UserRepository userRepository;

    @RequestMapping(value = "/", method = RequestMethod.GET)
    public List<User> getUsers() {
        List<User> response = new LinkedList<>();
        userRepository.findAll().forEach(response::add);
        return response;
    }

    @RequestMapping(value = "/{username}", method = RequestMethod.GET)
    public User getUser(HttpServletRequest request, @PathVariable("username") String username) {
        Object requestAttribute = request.getAttribute("claims");
        if (requestAttribute instanceof Claims) {
            Claims claims = (Claims) requestAttribute;
            String claimUsername = (String) claims.get("username");
            if (claimUsername != null && !username.equalsIgnoreCase(claimUsername)) {
                throw new AccessDeniedException("No access for requested entity");
            }
        }
        return userRepository.findOneByUsername(username);
    }
}
