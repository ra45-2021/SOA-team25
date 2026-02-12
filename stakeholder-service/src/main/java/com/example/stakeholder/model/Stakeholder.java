package com.example.stakeholder.model;

import jakarta.persistence.*;
import lombok.Data;

@Entity
@Table(name = "stakeholders")
@Data
public class Stakeholder {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "user_id", unique = true, nullable = false)
    private Long userId; 

    private String firstName;
    private String lastName;
    private String profilePicUrl;
    
    @Column(columnDefinition = "TEXT")
    private String biography;
    
    private String motto;
}