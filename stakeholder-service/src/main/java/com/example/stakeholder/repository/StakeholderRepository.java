package com.example.stakeholder.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.example.stakeholder.model.Stakeholder;

import java.util.Optional;

public interface StakeholderRepository extends JpaRepository<Stakeholder, Long> {
    Optional<Stakeholder> findByUserId(Long userId);
}