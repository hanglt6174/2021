package com.springexa.timeline.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.servlet.ModelAndView;

import com.springexa.timeline.model.User;
import com.springexa.timeline.service.UserService;

@Controller
public class IndexController {

    @Autowired
    private UserService userService;

    @GetMapping(value="/index")
    public ModelAndView home(){
        ModelAndView modelAndView = new ModelAndView();
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        User user = userService.findUserByEmail(auth.getName());
        modelAndView.addObject("user", "Welcome  (" + user.getEmail() + ")");
        modelAndView.addObject("message","index");
        modelAndView.setViewName("index");
        return modelAndView;
    }


}
