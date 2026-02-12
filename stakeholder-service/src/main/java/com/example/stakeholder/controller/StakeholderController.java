package com.example.stakeholder.controller;

import com.example.stakeholder.model.Stakeholder;
import com.example.stakeholder.repository.StakeholderRepository;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/stakeholders")
public class StakeholderController {

    private final StakeholderRepository repository;

    public StakeholderController(StakeholderRepository repository) {
        this.repository = repository;
    }

    @GetMapping
    public List<Stakeholder> getAll() {
        return repository.findAll();
    }

    @GetMapping("/user/{userId}")
    public Stakeholder getByUserId(@PathVariable Long userId) {
        return repository.findByUserId(userId).orElse(new Stakeholder());
    }

    @PostMapping
    public Stakeholder save(@RequestBody Stakeholder stakeholder) {
        return repository.save(stakeholder);
    }
}